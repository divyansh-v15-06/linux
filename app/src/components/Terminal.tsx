import React, { useEffect, useRef } from 'react';
import { Terminal as Xterm } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { parseCommand } from '../utils/parser';
import '@xterm/xterm/css/xterm.css';

export const Terminal: React.FC = () => {
  const containerRef = useRef<HTMLDivElement>(null);
  const termRef = useRef<Xterm | null>(null);

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

    // Use a small timeout to let container layout settle and fonts load
    const fitTimeout = setTimeout(() => {
      fitAddon.fit();
    }, 150);

    // Keep track of shell state
    let currentPath = ['home', 'guest'];
    let inputBuffer = '';
    const commandHistory: string[] = [];
    let historyIndex = -1;

    // Helper to format prompt
    const getPrompt = (path: string[]) => {
      let displayPath = '';
      if (path.length >= 2 && path[0] === 'home' && path[1] === 'guest') {
        const remaining = path.slice(2);
        displayPath = '~' + (remaining.length > 0 ? '/' + remaining.join('/') : '');
      } else {
        displayPath = '/' + path.join('/');
      }
      return `\x1b[1;32mguest@linuxquest\x1b[0m:\x1b[1;34m${displayPath}\x1b[0m$ `;
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
    term.write(getPrompt(currentPath));

    // Handle keypresses
    term.onData((data) => {
      const code = data.charCodeAt(0);

      // Enter key
      if (code === 13) {
        term.writeln('');
        const trimmedCmd = inputBuffer.trim();
        if (trimmedCmd) {
          commandHistory.push(trimmedCmd);
        }
        historyIndex = -1;

        // Parse and execute
        const result = parseCommand(inputBuffer, currentPath, commandHistory);
        if (result.clear) {
          term.clear();
        } else if (result.output) {
          term.writeln(result.output);
        }

        if (result.newPath) {
          currentPath = result.newPath;
        }

        inputBuffer = '';
        term.write(getPrompt(currentPath));
      }
      // Backspace
      else if (code === 127) {
        if (inputBuffer.length > 0) {
          inputBuffer = inputBuffer.slice(0, -1);
          term.write('\b \b');
        }
      }
      // ANSI escape sequences (Arrow keys, etc.)
      else if (code === 27) {
        if (data.startsWith('\x1b[A')) { // Up arrow
          if (commandHistory.length > 0) {
            if (historyIndex === -1) {
              historyIndex = commandHistory.length - 1;
            } else if (historyIndex > 0) {
              historyIndex--;
            }
            // Clear current input from terminal
            for (let i = 0; i < inputBuffer.length; i++) {
              term.write('\b \b');
            }
            inputBuffer = commandHistory[historyIndex];
            term.write(inputBuffer);
          }
        } else if (data.startsWith('\x1b[B')) { // Down arrow
          if (historyIndex !== -1) {
            // Clear current input
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
      // Ctrl+V or Command+V (macOS) or Ctrl+Shift+V
      const isPasteCombo = (e.ctrlKey || e.metaKey) && (e.key === 'v' || e.key === 'V');
      if (isPasteCombo) {
        term.writeln('\r\n\x1b[1;31m[SHIVA] lol. that\'s why AI is taking ur job.\x1b[0m');
        term.write(getPrompt(currentPath) + inputBuffer);
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
      term.write(getPrompt(currentPath) + inputBuffer);
    };

    if (termEl) {
      termEl.addEventListener('contextmenu', blockContextMenu);
      termEl.addEventListener('paste', handlePaste);
    }

    // Window resize handler
    const handleResize = () => {
      fitAddon.fit();
    };
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
  }, []);

  return (
    <div style={{ width: '100vw', height: '100vh', backgroundColor: '#0d1117', overflow: 'hidden' }}>
      <div 
        ref={containerRef} 
        style={{ 
          width: '100%', 
          height: '100%', 
          padding: '10px', 
          boxSizing: 'border-box' 
        }} 
      />
    </div>
  );
};
