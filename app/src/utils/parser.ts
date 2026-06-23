import { initialFS, getNodeByPath } from './fs';
import type { FSFile } from './fs';

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

export function parseCommand(
  rawLine: string,
  currentPath: string[],
  history: string[]
): ParseResult {
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
        // ls on a file just returns its name
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

      return { output: '', newPath: resolved };
    }

    case 'cat': {
      if (args.length === 0) {
        return { output: 'cat: missing file operand' };
      }
      const targetPath = args[0];
      const resolved = resolvePath(currentPath, targetPath);

      if (!resolved) {
        return { output: `cat: ${targetPath}: No such file or directory` };
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
        return { output: 'Usage: submit ANTARIKSHA{flag_value}' };
      }
      const flag = args[0];
      if (flag === 'ANTARIKSHA{TERMINAL_WAKES_2026}') {
        return { output: '\x1b[1;32m[CORRECT] Flag matches! Chapter 0 complete. +500 XP earned.\x1b[0m' };
      }
      return { output: '\x1b[1;31m[WRONG] Invalid flag. Check coordinates.txt or check spelling.\x1b[0m' };
    }

    case 'hint': {
      return { output: 'Hint: Read hint_coords.txt in your home directory (`cat hint_coords.txt`).' };
    }

    case 'start': {
      // Check if we are inside a chapter directory
      const isCh0 = currentPath.join('/').endsWith('antariksha/ch0');
      const isCh1 = currentPath.join('/').endsWith('antariksha/ch1');

      if (isCh0) {
        return { output: 'CheerpX: Starting virtual environment for Bootcamp...\r\n[Note: CheerpX sandbox connection will boot in the next phase!]' };
      }
      if (isCh1) {
        return { output: 'CheerpX: Chapter 1 is currently LOCKED. Complete Chapter 0 Bootcamp first.' };
      }
      return { output: 'start: No active sandbox configuration for this directory.' };
    }

    default:
      return { output: `-bash: ${cmd}: command not found` };
  }
}
