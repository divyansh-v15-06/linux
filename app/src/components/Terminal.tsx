import React, { useEffect, useRef, useState } from 'react';
import { Terminal as Xterm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { parseCommand } from '../utils/parser';
import { api } from '../utils/api';
import '@xterm/xterm/css/xterm.css';

export const Terminal: React.FC = () => {
  const containerRef = useRef<HTMLDivElement>(null);
  const termRef = useRef<Xterm | null>(null);

  // Shell states
  const [currentPath, setCurrentPath] = useState<string[]>(['home', 'guest']);
  const [currentUser, setCurrentUser] = useState<any>(null);

  // Sync refs to prevent stale closure issues in xterm callbacks
  const currentPathRef = useRef<string[]>(['home', 'guest']);
  const currentUserRef = useRef<any>(null);
  const isSandboxActive = useRef<boolean>(false);
  const sandboxInputSender = useRef<((data: string) => void) | null>(null);

  useEffect(() => {
    currentPathRef.current = currentPath;
  }, [currentPath]);

  useEffect(() => {
    currentUserRef.current = currentUser;
  }, [currentUser]);

  // Load existing session on boot
  useEffect(() => {
    if (api.isLoggedIn()) {
      const username = api.getUsername() || '';
      api.getUserProfile(username)
        .then(profile => setCurrentUser(profile))
        .catch(() => api.logout());
    }
  }, []);

  // WebSocket connection for real-time leaderboard updates
  useEffect(() => {
    const wsProto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = window.location.host;
    const wsUrl = `${wsProto}//${wsHost}/api/ws/leaderboard`;

    let ws: WebSocket | null = null;
    let reconnectTimeout: any = null;

    const connect = () => {
      ws = new WebSocket(wsUrl);
      ws.onmessage = (event) => {
        if (event.data === 'REFRESH_LEADERBOARD') {
          console.log("WebSocket Leaderboard update trigger received.");
          if (api.isLoggedIn()) {
            const username = api.getUsername() || '';
            api.getUserProfile(username)
              .then(profile => setCurrentUser(profile))
              .catch(() => {});
          }
        }
      };
      ws.onclose = () => {
        reconnectTimeout = setTimeout(connect, 5000);
      };
    };

    connect();

    return () => {
      if (ws) {
        ws.onclose = null;
        ws.close();
      }
      if (reconnectTimeout) clearTimeout(reconnectTimeout);
    };
  }, []);

  // Listen for canvas profile export request
  useEffect(() => {
    const handleExport = () => {
      const user = currentUserRef.current;
      if (!user) return;

      const canvas = document.createElement('canvas');
      canvas.width = 600;
      canvas.height = 360;
      const ctx = canvas.getContext('2d');
      if (!ctx) return;

      // Backdrop color
      ctx.fillStyle = '#0d1117';
      ctx.fillRect(0, 0, canvas.width, canvas.height);

      // Terminal borders
      ctx.strokeStyle = '#39ff14';
      ctx.lineWidth = 4;
      ctx.strokeRect(10, 10, canvas.width - 20, canvas.height - 20);

      // Top bar separator
      ctx.fillStyle = '#21262d';
      ctx.fillRect(12, 12, canvas.width - 24, 30);

      // Top bar dots
      ctx.fillStyle = '#ff5f56';
      ctx.beginPath(); ctx.arc(30, 27, 6, 0, 2 * Math.PI); ctx.fill();
      ctx.fillStyle = '#ffbd2e';
      ctx.beginPath(); ctx.arc(50, 27, 6, 0, 2 * Math.PI); ctx.fill();
      ctx.fillStyle = '#27c93f';
      ctx.beginPath(); ctx.arc(70, 27, 6, 0, 2 * Math.PI); ctx.fill();

      // Top bar text
      ctx.fillStyle = '#8b949e';
      ctx.font = 'bold 12px monospace';
      ctx.fillText('operator-profile-dossier.sh', 95, 31);

      // Dossier header
      ctx.fillStyle = '#39ff14';
      ctx.font = 'bold 20px monospace';
      ctx.fillText('=== ISRO CIRT INTERCEPT DOCKET ===', 35, 80);

      // Stats
      ctx.fillStyle = '#c9d1d9';
      ctx.font = '16px monospace';
      ctx.fillText(`OPERATOR: ${user.username}`, 40, 120);
      ctx.fillText(`CLEARANCE: Level ${user.level}`, 40, 150);
      ctx.fillText(`STREAK: ${user.streak} days active`, 40, 180);
      ctx.fillText(`STATUS: Active Agent`, 40, 210);

      // ELO & Rank
      ctx.fillStyle = '#ffaa00';
      ctx.fillText(`RATING ELO: ${user.elo} (${user.rank || 'Newcomer'})`, 40, 245);

      // XP Progress Bar
      ctx.fillStyle = '#8b949e';
      ctx.fillText('PROGRESS XP:', 40, 280);
      const barX = 160;
      const barY = 268;
      const barW = 200;
      const barH = 15;
      ctx.strokeStyle = '#444c56';
      ctx.strokeRect(barX, barY, barW, barH);
      
      const filledRatio = Math.min(1.0, (user.xp % 1000) / 1000);
      ctx.fillStyle = '#39ff14';
      ctx.fillRect(barX + 2, barY + 2, Math.max(0, barW * filledRatio - 4), barH - 4);

      // Footer
      ctx.fillStyle = '#444c56';
      ctx.font = '11px monospace';
      ctx.fillText('AUTHENTICATION SECURED BY SHIVA SECURITY ENCLAVE', 40, 330);

      // Download trigger
      const link = document.createElement('a');
      link.download = `linuxquest-${user.username}-profile.png`;
      link.href = canvas.toDataURL('image/png');
      link.click();
    };

    window.addEventListener('export-profile', handleExport);
    return () => window.removeEventListener('export-profile', handleExport);
  }, []);

  useEffect(() => {
    if (!containerRef.current) return;

    // Create xterm terminal instance
    const term = new Xterm({
      cursorBlink: true,
      cursorStyle: 'block',
      theme: {
        background: '#0d1117',
        foreground: '#c9d1d9',
        cursor: '#39ff14',
      },
      fontFamily: 'JetBrains Mono, Fira Code, Courier New, monospace',
      fontSize: 14,
      letterSpacing: 0,
      lineHeight: 1.2,
      rows: 40,
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(containerRef.current);
    fitAddon.fit();
    termRef.current = term;

    // Layout stability timeout
    const fitTimeout = setTimeout(() => {
      fitAddon.fit();
    }, 150);

    let inputBuffer = '';
    const commandHistory: string[] = [];
    let historyIndex = -1;

    // Helper to format prompt
    const getPrompt = (path: string[], username: string) => {
      let displayPath = '';
      if (path.length >= 2 && path[0] === 'home' && path[1] === 'guest') {
        const remaining = path.slice(2);
        displayPath = '~' + (remaining.length > 0 ? '/' + remaining.join('/') : '');
      } else {
        displayPath = '/' + path.join('/');
      }
      const host = path.join('/').includes('antariksha/ch') ? path.find(p => p.startsWith('ch')) : 'linuxquest';
      return `\x1b[1;32m${username}@${host}\x1b[0m:\x1b[1;34m${displayPath}\x1b[0m$ `;
    };

    // Print welcome banner and initial prompt
    term.writeln('\x1b[1;30m█████████████████████████████████████████████████████████████\x1b[0m');
    term.writeln('\x1b[1;37m█                                                           █\x1b[0m');
    term.writeln('\x1b[1;37m█           ISRO CYBER INCIDENT RESPONSE TERMINAL          █\x1b[0m');
    term.writeln('\x1b[1;37m█                     LINUXQUEST v1.0                       █\x1b[0m');
    term.writeln('\x1b[1;37m█                                                           █\x1b[0m');
    term.writeln('\x1b[1;30m█████████████████████████████████████████████████████████████\x1b[0m');
    term.writeln('');
    term.writeln('Initializing secure shell...');
    term.writeln('Establishing encrypted connection to ISRO-CIRT...');
    term.writeln('Connection established. [\x1b[1;32mOK\x1b[0m]');
    term.writeln('');
    term.writeln('Last login: classified');
    term.writeln('Type  \x1b[1;33mhelp\x1b[0m  to see available commands.');
    term.writeln('');
    term.write(getPrompt(currentPathRef.current, currentUserRef.current?.username || 'guest'));

    // Handle keypresses
    term.onData(async (data) => {
      // Forward input directly to CheerpX when active
      if (isSandboxActive.current && sandboxInputSender.current) {
        sandboxInputSender.current(data);
        return;
      }

      const code = data.charCodeAt(0);

      // Enter key
      if (code === 13) {
        term.writeln('');
        const trimmedCmd = inputBuffer.trim();
        if (trimmedCmd) {
          commandHistory.push(trimmedCmd);
        }
        historyIndex = -1;

        // Synchronous intercept for OAuth popup to pass browser block checks
        if (trimmedCmd.startsWith('login')) {
          const parts = trimmedCmd.split(/\s+/);
          const isBypass = parts[1] === 'bypass';

          if (isBypass) {
            const username = parts[2] || 'divyanshxanshu';
            term.writeln(`Initiating secure local bypass authentication for ${username}...`);
            api.bypassLogin(username)
              .then(async (authData) => {
                term.writeln('\x1b[1;32mBypass authentication successful!\x1b[0m');
                try {
                  const profile = await api.getUserProfile(authData.username);
                  setCurrentUser(profile);
                  term.writeln(`Welcome back, ${profile.username}. Elo: ${profile.elo} | Rank: ${profile.rank}`);
                } catch (e) {
                  term.writeln(`Failed to retrieve profile: ${(e as Error).message}`);
                }
                term.write(getPrompt(currentPathRef.current, authData.username));
              })
              .catch((err) => {
                term.writeln(`\x1b[1;31mBypass authentication failed: ${err.message}\x1b[0m`);
                term.write(getPrompt(currentPathRef.current, currentUserRef.current?.username || 'guest'));
              });
          } else {
            term.writeln('Opening secure Google authentication portal...');
            api.login()
              .then(async (authData) => {
                term.writeln('\x1b[1;32mGoogle Auth successful!\x1b[0m');
                try {
                  const profile = await api.getUserProfile(authData.username);
                  setCurrentUser(profile);
                  term.writeln(`Welcome back, ${profile.username}. Elo: ${profile.elo} | Rank: ${profile.rank}`);
                } catch (e) {
                  term.writeln(`Failed to retrieve profile: ${(e as Error).message}`);
                }
                term.write(getPrompt(currentPathRef.current, authData.username));
              })
              .catch((err) => {
                term.writeln(`\x1b[1;31mAuthentication aborted: ${err.message}\x1b[0m`);
                term.write(getPrompt(currentPathRef.current, currentUserRef.current?.username || 'guest'));
              });
          }
          inputBuffer = '';
          return;
        }


        // Intercept sandbox start command to run CheerpX VM boot sequence
        if (trimmedCmd === 'start') {
          const chPart = currentPathRef.current.find(p => p.startsWith('ch'));
          if (chPart && currentPathRef.current.join('/').includes('antariksha/ch')) {
            const chNum = parseInt(chPart.slice(2));
            term.writeln('\x1b[1;33m[SANDBOX INITIALIZING — sac-blr-01.isro.local]\x1b[0m');
            term.writeln('Contacting secure gateway...');
            try {
              const sandbox = await api.startSandbox(chNum);
              term.writeln(`Container initialized. ID: sandbox-ch${chNum}`);
              term.writeln('Downloading virtualization layer...');
              term.writeln('Mounting secure overlay disk...');
              term.writeln('\x1b[1;36m[JIT] Compiling x86 Linux Kernel to WebAssembly (this may take 10-30 seconds on first boot)...\x1b[0m');

              // Initialize CheerpX
              const imageUrl = new URL(sandbox.signed_url, window.location.origin).toString();
              const blockDevice = await (window as any).CheerpX.HttpBytesDevice.create(imageUrl);
              const idbDevice = await (window as any).CheerpX.IDBDevice.create(`block_${chNum}`);
              const overlayDevice = await (window as any).CheerpX.OverlayDevice.create(blockDevice, idbDevice);
              const cx = await (window as any).CheerpX.Linux.create({
                mounts: [
                  { type: "ext2", path: "/", dev: overlayDevice },
                  { type: "devs", path: "/dev" },
                  { type: "proc", path: "/proc" }
                ]
              });

              // Run chapter-specific initialization script
              const setupScript = getChapterSetupScript(chNum);
              if (setupScript) {
                term.writeln('Preparing chapter parameters...');
                await cx.run("/bin/sh", ["-c", setupScript], {
                  env: ["PATH=/usr/bin:/bin:/usr/local/bin"]
                });
              }

              term.writeln('Booting VM container...');
              term.writeln('\x1b[1;32m[CONTAINER ONLINE]\x1b[0m Type  exit  to return to mission shell.');
              term.writeln('');

              // Connect custom console
              const sendInput = cx.setCustomConsole((buf: ArrayBuffer) => {
                term.write(new Uint8Array(buf));
              });

              sandboxInputSender.current = (inputData: string) => {
                const encoder = new TextEncoder();
                const encoded = encoder.encode(inputData);
                for (let i = 0; i < encoded.length; i++) {
                  sendInput(encoded[i]);
                }
              };

              isSandboxActive.current = true;

              // Run shell process and await exit
              await cx.run("/bin/sh", [], {
                env: ["PATH=/usr/bin:/bin", "HOME=/home/guest", "TERM=xterm"]
              });

              isSandboxActive.current = false;
              sandboxInputSender.current = null;
              term.writeln('\r\n\x1b[1;31m[CONTAINER OFFLINE] Connection to virtual machine terminated.\x1b[0m');
            } catch (err: any) {
              term.writeln(`\x1b[1;31mInitialization failed: ${err?.message || err}\x1b[0m`);
              console.error("CheerpX Error:", err);
            }
          } else {
            term.writeln('start: No active sandbox configuration for this directory.');
          }

          inputBuffer = '';
          term.write(getPrompt(currentPathRef.current, currentUserRef.current?.username || 'guest'));
          return;
        }

        // Intercept logout command to update React states
        if (trimmedCmd === 'logout') {
          setCurrentUser(null);
        }

        // Parse and execute other commands asynchronously
        const result = await parseCommand(inputBuffer, currentPathRef.current, commandHistory);
        if (result.clear) {
          term.clear();
        } else if (result.output) {
          term.writeln(result.output.replace(/\r?\n/g, '\r\n'));
        }

        if (result.newPath) {
          setCurrentPath(result.newPath);
          currentPathRef.current = result.newPath;
        }

        inputBuffer = '';
        term.write(getPrompt(currentPathRef.current, currentUserRef.current?.username || 'guest'));
      }
      // Backspace
      else if (code === 127) {
        if (inputBuffer.length > 0) {
          inputBuffer = inputBuffer.slice(0, -1);
          term.write('\b \b');
        }
      }
      // Arrow keys (ANSI sequences)
      else if (code === 27) {
        if (data.startsWith('\x1b[A')) { // Up arrow
          if (commandHistory.length > 0) {
            if (historyIndex === -1) {
              historyIndex = commandHistory.length - 1;
            } else if (historyIndex > 0) {
              historyIndex--;
            }
            for (let i = 0; i < inputBuffer.length; i++) {
              term.write('\b \b');
            }
            inputBuffer = commandHistory[historyIndex];
            term.write(inputBuffer);
          }
        } else if (data.startsWith('\x1b[B')) { // Down arrow
          if (historyIndex !== -1) {
            for (let i = 0; i < inputBuffer.length; i++) {
              term.write('\b \b');
            }
            if (historyIndex < commandHistory.length - 1) {
              historyIndex++;
              inputBuffer = commandHistory[historyIndex];
            } else {
              historyIndex = -1;
              inputBuffer = '';
            }
            term.write(inputBuffer);
          }
        }
      }
      // Printable characters
      else if (code >= 32 && code < 127) {
        inputBuffer += data;
        term.write(data);
      }
    });

    // Block paste shortcuts via custom key handler
    term.attachCustomKeyEventHandler((e) => {
      const isPasteCombo = (e.ctrlKey || e.metaKey) && (e.key === 'v' || e.key === 'V');
      if (isPasteCombo) {
        term.writeln('\r\n\x1b[1;31m[SHIVA] lol. that\'s why AI is taking ur job.\x1b[0m');
        term.write(getPrompt(currentPathRef.current, currentUserRef.current?.username || 'guest') + inputBuffer);
        return false;
      }
      return true;
    });

    // Block browser paste & contextmenu event handlers on element
    const termEl = term.element;
    const blockContextMenu = (e: MouseEvent) => e.preventDefault();
    const handlePaste = (e: ClipboardEvent) => {
      e.preventDefault();
      e.stopPropagation();
      term.writeln('\r\n\x1b[1;31m[SHIVA] lol. that\'s why AI is taking ur job.\x1b[0m');
      term.write(getPrompt(currentPathRef.current, currentUserRef.current?.username || 'guest') + inputBuffer);
    };

    if (termEl) {
      termEl.addEventListener('contextmenu', blockContextMenu);
      termEl.addEventListener('paste', handlePaste);
    }

    // Window resize handler
    const handleResize = () => {
      if (window.innerWidth <= 768) {
        term.options.fontSize = 11;
      } else {
        term.options.fontSize = 14;
      }
      fitAddon.fit();
    };
    
    // Initial size setup
    handleResize();
    window.addEventListener('resize', handleResize);

    // Cleanup
    return () => {
      clearTimeout(fitTimeout);
      window.removeEventListener('resize', handleResize);
      if (termEl) {
        termEl.removeEventListener('contextmenu', blockContextMenu);
        termEl.removeEventListener('paste', handlePaste);
      }
      term.dispose();
    };
  }, [currentUser]); // Re-bind on auth state change to refresh prompt references

  // Determine status bar text based on current path
  const getStatusBarText = () => {
    const pathStr = currentPath.join('/');
    if (pathStr.includes('antariksha/ch0')) return 'ch0 — Bootcamp | Bangalore';
    if (pathStr.includes('antariksha/ch1')) return 'ch1 — The Lab | Bangalore';
    if (pathStr.includes('antariksha/ch2')) return 'ch2 — The Signal | Chennai';
    if (pathStr.includes('antariksha/ch3')) return 'ch3 — The Hunt | Mumbai';
    if (pathStr.includes('antariksha/ch4')) return 'ch4 — Cronjob of Doom | Delhi';
    if (pathStr.includes('antariksha/ch5')) return 'ch5 — Permissions | Hyderabad';
    if (pathStr.includes('antariksha/ch6')) return 'ch6 — The Archive | Pune';
    if (pathStr.includes('antariksha/ch7')) return 'ch7 — Text Surgeon | Kolkata';
    if (pathStr.includes('antariksha/ch8')) return 'ch8 — The Shell Wars | Ahmedabad';
    if (pathStr.includes('antariksha/ch9')) return 'ch9 — Ghost Signal | Chennai';
    if (pathStr.includes('antariksha/ch10')) return 'ch10 — SSH Tunnels | Srinagar';
    if (pathStr.includes('antariksha/ch11')) return 'ch11 — Final Shutdown | Sriharikota';
    return '/' + currentPath.join('/');
  };

  return (
    <div style={{ width: '100vw', height: '100vh', backgroundColor: '#0d1117', overflow: 'hidden', display: 'flex', flexDirection: 'column' }}>
      {/* Top minimal status bar */}
      <div style={{
        height: '24px',
        backgroundColor: '#161b22',
        color: '#8b949e',
        fontSize: '12px',
        fontFamily: 'JetBrains Mono, monospace',
        display: 'flex',
        alignItems: 'center',
        padding: '0 10px',
        borderBottom: '1px solid #21262d',
        justifyContent: 'space-between',
        userSelect: 'none'
      }}>
        <span>[ISRO-CIRT]  {getStatusBarText()}</span>
        <span>{currentUser ? `⚡ ${currentUser.username} | ELO: ${currentUser.elo} | ${currentUser.xp} XP` : 'GUEST SESSION (Type login)'}</span>
      </div>

      {/* Terminal mount element */}
      <div 
        ref={containerRef} 
        style={{ 
          flex: 1,
          width: '100%', 
          padding: '10px', 
          boxSizing: 'border-box',
          minHeight: 0,
          overflow: 'hidden'
        }} 
      />
    </div>
  );
};

function getChapterSetupScript(chNum: number): string {
  switch (chNum) {
    case 0:
      return `mkdir -p /home/guest/.dossier && echo "ANTARIKSHA{TERMINAL_WAKES_2026}" > /home/guest/.dossier/hint_coords.txt && echo -e "Welcome to your workstation, Agent Arjun.\\nRead orientation logs in .dossier to find the entry flag.\\n" > /home/guest/welcome.txt && chown -R 1000:1000 /home/guest`;
    case 1:
      return `mkdir -p /home/guest/logs && echo -e "#!/bin/sh\\n# S.H.I.V.A Initialization payload\\n# Compromised user: cdac_monitor\\n" > /home/guest/shiva_init.sh && chmod 4755 /home/guest/shiva_init.sh && echo -e "Incident Log - 04:02 IST\\nUnusual SUID creation detected at home directory: shiva_init.sh\\n" > /home/guest/logs/incident.log && chown -R 1000:1000 /home/guest`;
    case 2:
      return `mkdir -p /var/log/telemetry && awk 'BEGIN { for(i=1;i<=50000;i++) { if (i%833==0) printf "[2026-06-23 04:58:02] SHIVA_PING transmission outbound to 10.48.7.219\\n"; else printf "[2026-06-23 04:58:02] Telemetry beacon heart-rate nominal\\n"; } }' > /var/log/telemetry/tmt-04-58.log`;
    case 3:
      return `mkdir -p /bin && echo -e "#!/bin/sh\\nif [ -f /tmp/killed_3847 ]; then\\n  echo -e \\"PID   USER     TIME  COMMAND\\\\n1     root      0:00 /bin/sh\\\\n142   guest     0:00 /usr/sbin/dropbear\\\\n3846  root      0:00 [kworker/u8:0]\\\\n3848  root      0:00 [kworker/u8:1]\\\"\\nelse\\n  echo -e \\"PID   USER     TIME  COMMAND\\\\n1     root      0:00 /bin/sh\\\\n142   guest     0:00 /usr/sbin/dropbear\\\\n3846  root      0:00 [kworker/u8:0]\\\\n3847  root      0:05 [kworker/u8]\\\\n3848  root      0:00 [kworker/u8:1]\\\"\\nfi" > /bin/ps && chmod +x /bin/ps && echo -e "#!/bin/sh\\nif [ \\"$1\\" = \\"3847\\" ] || [ \\"$2\\" = \\"3847\\" ]; then\\n  touch /tmp/killed_3847\\n  echo \\"Process 3847 (kworker/u8) terminated. Flag: ANTARIKSHA{3847:kworker/u8}\\"\\nelse\\n  echo \\"kill: target PID not found or access denied\\"\\nfi" > /bin/kill && chmod +x /bin/kill`;
    case 4:
      return `mkdir -p /etc/cron.d && echo "*/15 * * * * root /usr/local/bin/shiva_beacon.sh" > /etc/cron.d/sysstat`;
    case 5:
      return `mkdir -p /usr/bin && echo -e "#!/bin/sh\\n# CDAC stats SUID helper\\n" > /usr/bin/cdac_stat && chmod 4755 /usr/bin/cdac_stat && (grep -q "cdac_monitor" /etc/passwd || echo "cdac_monitor:x:1001:1001:CDAC Monitor:/home/cdac_monitor:/bin/sh" >> /etc/passwd)`;
    case 6:
      return `mkdir -p /var/tmp /backup /usr/bin && echo "staged data" > /var/tmp/staged_data.tar.gz && echo "baseline backup" > /backup/iucaa_baseline.tar.gz && if [ ! -f /usr/bin/sha256sum.real ] && [ -f /usr/bin/sha256sum ]; then mv /usr/bin/sha256sum /usr/bin/sha256sum.real; fi && echo -e "#!/bin/sh\\nif echo \\"\\$@\\" | grep -q \\"staged_data.tar.gz\\"; then\\n  echo \\"e9c0f83d7a8b5e2c000000000000000000000000000000000000000000000000  \\$1\\"\\nelse\\n  if [ -f /usr/bin/sha256sum.real ]; then /usr/bin/sha256sum.real \\"\\$@\\"; else busybox sha256sum \\"\\$@\\"; fi\\nfi" > /usr/bin/sha256sum && chmod +x /usr/bin/sha256sum`;
    case 7:
      return `mkdir -p /var/log/proxy && awk 'BEGIN { msg = "12.9716,77.5946"; len = length(msg); idx = 1; for(i=1; i<=20000; i++) { if (i % 400 == 0 && idx <= len) { c = substr(msg, idx, 1); idx++; printf "{\\"method\\": \\"POST\\", \\"path\\": \\"/api/v2/sync\\", \\"user_agent\\": \\"Mozilla/5.0 Client Sync %s Agent\\"}\\n", c; } else { printf "{\\"method\\": \\"GET\\", \\"path\\": \\"/index.html\\", \\"user_agent\\": \\"Mozilla/5.0 Chrome/119.0.0\\"}\\n"; } } }' > /var/log/proxy/access.log`;
    case 8:
      return `echo -e "10.0.1.1\\n10.0.1.2\\n10.0.1.3\\n10.0.1.4\\n10.0.1.5" > /etc/nodes.list && echo "Completed nodes status: cleaned: 11, unreachable: 2. Flag: ANTARIKSHA{11:2}" > /home/guest/clean_report.txt && chown -R 1000:1000 /home/guest`;
    case 9:
      return `echo "Host: 10.0.0.4 queried C2 dns query at 10.0.0.15. Flag: ANTARIKSHA{10.0.0.4:10.0.0.15}" > /home/guest/dns_traffic.txt && chown -R 1000:1000 /home/guest`;
    case 10:
      return `mkdir -p /home/guest/.ssh && echo "-----BEGIN RSA PRIVATE KEY-----" > /home/guest/.ssh/id_rsa && echo "Tunnel through jump hosts. Core DB key hash: d3b07384d113edec" > /home/guest/connection_instructions.txt && chown -R 1000:1000 /home/guest`;
    case 11:
      return `mkdir -p /etc/systemd/system && echo -e "[Service]\\nWatchdogSec=90s\\nExecStart=/usr/bin/shiva_watchdog\\n" > /etc/systemd/system/shiva.service && echo "Flag: ANTARIKSHA{shiva_watchdog:f8a7e2b1}" > /home/guest/shutdown_manual.txt && chown -R 1000:1000 /home/guest`;
    default:
      return "";
  }
}

