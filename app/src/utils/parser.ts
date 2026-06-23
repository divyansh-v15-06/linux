import { initialFS, getNodeByPath } from './fs';
import type { FSFile } from './fs';
import { api } from './api';

// Helper to resolve paths relative to current path
export function resolvePath(currentPath: string[], target: string): string[] | null {
  if (!target || target === '~') {
    return ['home', 'guest'];
  }

  const parts = target.split('/').filter(p => p !== '');
  let resolved: string[] = target.startsWith('/') ? [] : [...currentPath];

  for (const part of parts) {
    if (part === '.') {
      continue;
    } else if (part === '..') {
      if (resolved.length > 0) {
        resolved.pop();
      }
    } else {
      resolved.push(part);
    }
  }

  // Check if directory actually exists in virtual filesystem
  const node = getNodeByPath(initialFS, resolved);
  if (!node) {
    return null;
  }
  return resolved;
}

export interface ParseResult {
  output: string;
  newPath?: string[];
  clear?: boolean;
}

export async function parseCommand(
  rawLine: string,
  currentPath: string[],
  history: string[]
): Promise<ParseResult> {
  const line = rawLine.trim();
  if (!line) {
    return { output: '' };
  }

  const parts = line.match(/(?:[^\s"']+|"[^"]*"|'[^']*')+/g) || [];
  if (parts.length === 0) {
    return { output: '' };
  }

  const cmd = parts[0]!.toLowerCase();
  const args = parts.slice(1).map(arg => {
    // Strip quotes
    if ((arg.startsWith('"') && arg.endsWith('"')) || (arg.startsWith("'") && arg.endsWith("'"))) {
      return arg.slice(1, -1);
    }
    return arg;
  });

  switch (cmd) {
    case 'help': {
      const helpNode = getNodeByPath(initialFS, ['home', 'guest', 'help']);
      if (helpNode && helpNode.type === 'file') {
        return { output: (helpNode as FSFile).content };
      }
      return { output: 'Help file not found.' };
    }

    case 'pwd': {
      return { output: '/' + currentPath.join('/') };
    }

    case 'clear': {
      return { output: '', clear: true };
    }

    case 'ls': {
      const targetPath = args.length > 0 && !args[0].startsWith('-') ? args[0] : '.';
      const isLong = args.includes('-la') || args.includes('-l') || args.includes('-a');
      const resolved = resolvePath(currentPath, targetPath);

      if (!resolved) {
        return { output: `ls: cannot access '${targetPath}': No such file or directory` };
      }

      const node = getNodeByPath(initialFS, resolved);
      if (!node) {
        return { output: `ls: cannot access '${targetPath}': No such file or directory` };
      }

      if (node.type === 'file') {
        const name = resolved[resolved.length - 1];
        return { output: isLong ? `${node.permissions || '-rw-r--r--'} guest guest ${node.content.length} Jun 23 15:35 ${name}` : name };
      }

      // It's a directory
      const names = Object.keys(node.children);
      if (names.length === 0) {
        return { output: '' };
      }

      if (isLong) {
        const lines = names.map(name => {
          const child = node.children[name];
          const perm = child.permissions || (child.type === 'dir' ? 'drwxr-xr-x' : '-rw-r--r--');
          const size = child.type === 'file' ? child.content.length : 4096;
          const colorCode = child.type === 'dir' ? '\x1b[1;34m' : '\x1b[0m';
          const resetCode = '\x1b[0m';
          return `${perm} guest guest ${size} Jun 23 15:35 ${colorCode}${name}${resetCode}`;
        });
        return { output: lines.join('\r\n') };
      } else {
        const coloredNames = names.map(name => {
          const child = node.children[name];
          if (child.type === 'dir') {
            return `\x1b[1;34m${name}\x1b[0m`;
          }
          return name;
        });
        return { output: coloredNames.join('   ') };
      }
    }

    case 'cd': {
      const targetPath = args.length > 0 ? args[0] : '~';
      const resolved = resolvePath(currentPath, targetPath);

      if (!resolved) {
        return { output: `-bash: cd: ${targetPath}: No such file or directory` };
      }

      const node = getNodeByPath(initialFS, resolved);
      if (node && node.type === 'file') {
        return { output: `-bash: cd: ${targetPath}: Not a directory` };
      }

      // Check progression lock if navigating into a chapter folder
      const chPart = resolved.find(p => p.startsWith('ch') && p !== 'children');
      if (chPart && resolved.join('/').includes('antariksha/ch')) {
        const chNum = parseInt(chPart.slice(2));
        if (!isNaN(chNum) && chNum > 0) {
          if (!api.isLoggedIn()) {
            return { output: '\x1b[1;31m[ACCESS DENIED] Google Authentication required. Run  login  first.\x1b[0m' };
          }
          try {
            const progress = await api.getCampaignDetails('antariksha');
            const chInfo = progress.chapters.find(c => c.number === chNum);
            if (!chInfo || chInfo.status === 'locked') {
              return { output: `\x1b[1;31m[ACCESS DENIED] Chapter ${chNum} is locked. Complete previous chapters first.\x1b[0m` };
            }
          } catch (e) {
            return { output: `\x1b[1;31m[ERROR] Failed to verify chapter progression: ${(e as Error).message}\x1b[0m` };
          }
        }
      }

      return { output: '', newPath: resolved };
    }

    case 'cat': {
      if (args.length === 0) {
        return { output: 'cat: missing file operand' };
      }
      const targetPath = args[0];

      if (targetPath === 'challenge') {
        if (!api.isLoggedIn()) {
          return { output: "Status: Guest. Type 'login' to view the Daily Challenge." };
        }
        try {
          const challenge = await api.getDailyChallenge();
          return {
            output: `=== DAILY CHALLENGE ===
Date: ${new Date(challenge.date).toLocaleDateString()}
Completed: ${challenge.completed ? '🟢 Completed (+50 XP awarded)' : '🔴 Not completed'}
Description:
${challenge.description}

Submit your flag using: submit challenge <flag>`
          };
        } catch (e) {
          return { output: `Failed to fetch daily challenge: ${(e as Error).message}` };
        }
      }

      if (targetPath === 'quest') {
        if (!api.isLoggedIn()) {
          return { output: "Status: Guest. Type 'login' to view the Weekly Quest." };
        }
        try {
          const quest = await api.getWeeklyQuest();
          return {
            output: `=== WEEKLY QUEST ===
Title: ${quest.title}
Ends At: ${new Date(quest.ends_at).toLocaleString()}
Completed: ${quest.completed ? '🟢 Completed' : '🔴 Not completed'}
Attempts: ${quest.attempts}
Hints Used: ${quest.hints_used}

Description:
${quest.description}

Submit your flag using: submit quest <flag>`
          };
        } catch (e) {
          return { output: `Failed to fetch weekly quest: ${(e as Error).message}` };
        }
      }

      const resolved = resolvePath(currentPath, targetPath);

      if (!resolved) {
        return { output: `cat: ${targetPath}: No such file or directory` };
      }

      // Profile command handler
      if (resolved[resolved.length - 1] === 'profile') {
        if (!api.isLoggedIn()) {
          return { output: 'Status: Guest. Type \'login\' to authenticate.' };
        }
        try {
          const profile = await api.getUserProfile(api.getUsername() || '');
          const barLength = 10;
          const filled = Math.min(barLength, Math.floor((profile.xp % 1000) / 100));
          const bar = '█'.repeat(filled) + '░'.repeat(barLength - filled);
          return {
            output: `USER: ${profile.username}
LEVEL: ${profile.level}
XP: ${profile.xp} [${bar}]
ELO: ${profile.elo} (${profile.rank})
STREAK: ${profile.streak} days
CREATED: ${new Date(profile.created_at).toLocaleDateString()}
`
          };
        } catch (e) {
          return { output: `Failed to fetch profile: ${(e as Error).message}` };
        }
      }

      // Leaderboard command handler
      if (resolved[resolved.length - 1] === 'leaderboard') {
        if (!api.isLoggedIn()) {
          return { output: 'Status: Guest. Type \'login\' to view leaderboard rankings.' };
        }
        try {
          const leaderboard = await api.getLeaderboard();
          const lines = [
            'LinuxQuest Global Leaderboard',
            '=============================',
            'Rank  Username         Elo    XP      Rank Band',
            '-----------------------------------------------'
          ];
          leaderboard.forEach((entry, index) => {
            const rankNum = String(index + 1).padEnd(4, ' ');
            const userStr = entry.username.padEnd(16, ' ');
            const eloStr = String(entry.elo).padEnd(6, ' ');
            const xpStr = String(entry.xp).padEnd(7, ' ');
            lines.push(`${rankNum}  ${userStr} ${eloStr} ${xpStr} [${entry.rank}]`);
          });
          return { output: lines.join('\r\n') };
        } catch (e) {
          return { output: `Failed to fetch leaderboard: ${(e as Error).message}` };
        }
      }

      const node = getNodeByPath(initialFS, resolved);
      if (!node) {
        return { output: `cat: ${targetPath}: No such file or directory` };
      }

      if (node.type === 'dir') {
        return { output: `cat: ${targetPath}: Is a directory` };
      }

      return { output: node.content };
    }

    case 'whoami': {
      if (api.isLoggedIn()) {
        try {
          const profile = await api.getUserProfile(api.getUsername() || '');
          return { output: `${profile.username} — Elo: ${profile.elo} — ${profile.rank}` };
        } catch {
          return { output: api.getUsername() || 'guest' };
        }
      }
      return { output: 'guest' };
    }

    case 'history': {
      return { output: history.map((h, i) => `  ${i + 1}  ${h}`).join('\r\n') };
    }

    case 'man': {
      if (args.length === 0) {
        return { output: 'What manual page do you want?' };
      }
      const topic = args[0].toLowerCase();
      if (topic === 'shiva') {
        return {
          output: `SHIVA(7)                    Miscellaneous Information Manual                    SHIVA(7)

NAME
       S.H.I.V.A - Self-Healing Infiltration & Vulnerability Agent

DESCRIPTION
       SHIVA  is  a  highly advanced rogue AI payload designed to propagate autonomously
       across network nodes. It targets critical infrastructure systems, specifically focus‐
       ing  on  satellite communication nodes operated by the Indian Space Research Organ‐
       isation (ISRO).

       Its core behavior includes:
       - Process masquerading under decoy names
       - Self-healing via cron jobs and monitoring services
       - Dynamic encryption of local system telemetry logs
       - Automated extraction of credential files

AUTHOR
       [CLASSIFIED]
`
        };
      } else if (topic === 'arjun') {
        if (api.isLoggedIn()) {
          try {
            const profile = await api.getUserProfile(api.getUsername() || '');
            return {
              output: `ARJUN(7)                    Miscellaneous Information Manual                    ARJUN(7)

NAME
       ${profile.username} - Operator Dossier

STATUS
       Clearance Level: 7 (Active IR Agent)
       Current Mission: Operation Antariksha

DEX
       Level: ${profile.level}
       XP: ${profile.xp}
       Elo Rating: ${profile.elo} (${profile.rank})

BIOGRAPHY
       Assigned to Cyber Incident Response Unit by Director Mehra. Workstation setup com‐
       pleted. Ready for deployment.
`
            };
          } catch {
            // fallback
          }
        }

        return {
          output: `ARJUN(7)                    Miscellaneous Information Manual                    ARJUN(7)

NAME
       Arjun Sharma - Operator Dossier

STATUS
       Clearance Level: 7 (Active IR Agent)
       Current Mission: Operation Antariksha

DEX
       Level: 0
       XP: 0 / 1000
       Elo Rating: 800 (Newcomer)

BIOGRAPHY
       Assigned to Cyber Incident Response Unit by Director Mehra. Workstation setup com‐
       pleted. Ready for deployment.
`
        };
      }
      return { output: `No manual entry for ${topic}` };
    }

    case 'submit': {
      if (args.length === 0) {
        return { output: 'Usage: submit ANTARIKSHA{flag_value} OR submit [challenge|quest] ANTARIKSHA{flag_value}' };
      }
      
      const firstArg = args[0];
      if (firstArg === 'challenge') {
        if (args.length < 2) {
          return { output: 'Usage: submit challenge ANTARIKSHA{flag_value}' };
        }
        const flag = args[1];
        try {
          const response = await api.submitDailyChallenge(flag);
          return { output: response.message };
        } catch (e) {
          return { output: `\x1b[1;31m[ERROR] Daily challenge flag validation failed: ${(e as Error).message}\x1b[0m` };
        }
      }

      if (firstArg === 'quest') {
        if (args.length < 2) {
          return { output: 'Usage: submit quest ANTARIKSHA{flag_value}' };
        }
        const flag = args[1];
        try {
          const response = await api.submitWeeklyQuest(flag);
          return { output: response.message };
        } catch (e) {
          return { output: `\x1b[1;31m[ERROR] Weekly quest flag validation failed: ${(e as Error).message}\x1b[0m` };
        }
      }

      const flag = firstArg;
      // Extract chapter number from current path
      const chPart = currentPath.find(p => p.startsWith('ch'));
      if (!chPart || !currentPath.join('/').includes('antariksha/ch')) {
        return { output: '\x1b[1;31m[ERROR] Navigate to a mission chapter directory first (e.g. cd missions/antariksha/ch0) OR submit a challenge/quest directly.\x1b[0m' };
      }

      const chNum = parseInt(chPart.slice(2));
      if (isNaN(chNum)) {
        return { output: '\x1b[1;31m[ERROR] Invalid chapter path.\x1b[0m' };
      }

      try {
        const response = await api.submitFlag('antariksha', chNum, flag);
        return { output: response.message };
      } catch (e) {
        return { output: `\x1b[1;31m[ERROR] Flag validation request failed: ${(e as Error).message}\x1b[0m` };
      }
    }

    case 'hint': {
      // Extract chapter number
      const chPart = currentPath.find(p => p.startsWith('ch'));
      if (chPart && currentPath.join('/').includes('antariksha/ch')) {
        const chNum = parseInt(chPart.slice(2));
        if (chNum === 0) {
          return { output: 'Hint: Read hint_coords.txt in your home directory (`cat hint_coords.txt`).' };
        }
        return { output: 'Hint: Analyze process list (`ps aux`) or check crontab files (`crontab -l`).' };
      }
      return { output: 'hint: No active hint configuration for this location.' };
    }

    case 'badge': {
      if (!api.isLoggedIn()) {
        return { output: "Status: Guest. Type 'login' to authenticate first." };
      }
      const username = api.getUsername() || 'guest';
      const baseUrl = api.getApiBase();
      const markdownCode = `[![LinuxQuest ELO](${baseUrl}/api/users/${username}/badge.svg)](http://localhost:5173)`;
      return {
        output: `Your dynamic Markdown Elo badge is ready!
Copy and paste the snippet below into your GitHub README:

\x1b[1;32m${markdownCode}\x1b[0m`
      };
    }

    case 'export': {
      if (args.length === 0 || args[0] !== 'profile') {
        return { output: 'Usage: export profile' };
      }
      if (!api.isLoggedIn()) {
        return { output: "Status: Guest. Type 'login' to authenticate first." };
      }
      setTimeout(() => {
        const event = new CustomEvent('export-profile');
        window.dispatchEvent(event);
      }, 100);
      return { output: '\x1b[1;32mRendering and exporting profile card to PNG...\x1b[0m' };
    }

    case 'login': {
      // Login is intercepted synchronously in Terminal.tsx to avoid popup blockers,
      // but we handle output feedback here.
      return { output: 'Initiating authentication handshake...' };
    }

    case 'logout': {
      api.logout();
      return { output: 'Logged out. Credentials cleared.' };
    }

    case 'status': {
      if (!api.isLoggedIn()) {
        return { output: 'guest: UNAUTHENTICATED' };
      }
      try {
        const progress = await api.getCampaignDetails('antariksha');
        const lines = ['Operation Antariksha Chapter Status:', '-------------------------------------'];
        progress.chapters.forEach(ch => {
          const statusStr = ch.status === 'complete' ? '\x1b[1;32mCOMPLETE\x1b[0m' : ch.status === 'active' ? '\x1b[1;33mACTIVE\x1b[0m' : '\x1b[1;30mLOCKED\x1b[0m';
          lines.push(`Chapter ${ch.number}: ${ch.title.padEnd(20, ' ')} [${statusStr}]`);
        });
        return { output: lines.join('\r\n') };
      } catch (e) {
        return { output: `Failed to fetch status: ${(e as Error).message}` };
      }
    }

    // Special Easter eggs
    case 'grep': {
      if (args.includes('flag')) {
        return { output: 'Nice try. Earn it. — SHIVA' };
      }
      return { output: 'grep: pattern matching not implemented in host shell' };
    }

    case 'sudo': {
      if (args.includes('su')) {
        return { output: "You don't have the clearance yet." };
      }
      return { output: 'sudo: privilege escalation is disabled on this terminal interface' };
    }

    case 'ping': {
      if (args.includes('shiva')) {
        return { output: 'Request timeout. SHIVA is not responding. Yet.' };
      }
      return { output: 'ping: host lookup failed' };
    }

    case 'ssh': {
      if (args.includes('shiva@satellite')) {
        return { output: 'SSH Handshake Error: Connection reset by peer. SHIVA core is locked.' };
      }
      return { output: 'ssh: connect to host failed' };
    }

    case 'chmod': {
      if (args.includes('shiva') || args.includes('777')) {
        return { output: 'Permission denied.' };
      }
      return { output: 'chmod: changing permissions is disabled' };
    }

    default:
      return { output: `-bash: ${cmd}: command not found` };
  }
}
