# LinuxQuest — Frontend Design Spec

> **For AI Agents:** The entire frontend is a terminal emulator. There is no traditional UI.
> No nav bars. No sidebars. No buttons. No modals. Just a shell.

---

## Core Principle

Every interaction on LinuxQuest happens through shell commands.
The "app" is a virtual filesystem. Navigation = Linux commands.

The player never leaves the terminal — not for menus, not for their profile,
not for leaderboards, not for mission briefs. Everything renders inside the shell.

**Copy-paste is disabled. Always. No exceptions.**

> *"Copy-paste is disabled. This is not a bug. LinuxQuest teaches commands through
> repetition. Typing `grep -r "SHIVA" /var/log/` ten times is the point.
> Your fingers will remember it when your terminal at work doesn't have hints."*

When a player tries to paste:
```
[PASTE DISABLED] Type the command. No shortcuts here.
```

---

## Boot Sequence

When a player first opens the site, they see this — nothing else:

```
█████████████████████████████████████████████████████████████
█                                                           █
█           ISRO CYBER INCIDENT RESPONSE TERMINAL          █
█                     LINUXQUEST v1.0                       █
█                                                           █
█████████████████████████████████████████████████████████████

Initializing secure shell...
Establishing encrypted connection to ISRO-CIRT...
Connection established. [OK]

Last login: classified

Type  help  to see available commands.

guest@linuxquest:~$ _
```

- Typing animation on the boot lines (20–40ms per character)
- Blinking block cursor `█`
- No buttons. No skip. Just the prompt.

---

## Virtual Filesystem Structure

The app's navigation tree maps 1:1 to a fake Linux filesystem.

```
/
└── home/
    └── guest/             ← starting directory (or arjun/ after login)
        ├── missions/
        │   ├── README         ← cat missions/README → campaign list
        │   └── antariksha/
        │       ├── README     ← cat antariksha/README → campaign overview
        │       ├── ch0/       ← cd ch0 → Bootcamp
        │       ├── ch1/       ← cd ch1 → The Lab
        │       ├── ch2/       ← cd ch2 → The Signal
        │       ├── ch3/       ← cd ch3 → The Hunt
        │       ├── ch4/       ← cd ch4 → Cronjob of Doom
        │       ├── ch5/       ← cd ch5 → Permissions
        │       ├── ch6/       ← cd ch6 → The Archive
        │       ├── ch7/       ← cd ch7 → Text Surgeon
        │       ├── ch8/       ← cd ch8 → The Shell Wars
        │       ├── ch9/       ← cd ch9 → Ghost Signal
        │       ├── ch10/      ← cd ch10 → SSH Tunnels
        │       └── ch11/      ← cd ch11 → Final Shutdown
        ├── leaderboard        ← cat leaderboard
        ├── profile            ← cat profile
        ├── daily              ← cat daily
        ├── tracks/
        │   ├── README
        │   ├── scripting/
        │   ├── vim/
        │   ├── git/
        │   ├── docker/
        │   └── kubernetes/
        └── help               ← cat help
```

---

## Command Reference

### Navigation

| Command | Output |
|---------|--------|
| `ls` | Lists files/dirs in current location |
| `ls -la` | Lists with metadata (locked chapters show as `----------`) |
| `cd missions` | Enters missions directory |
| `cd antariksha` | Enters Operation Antariksha |
| `cd ch1` | Enters Chapter 1 — displays mission brief, starts sandbox |
| `cd ..` | Goes up one level |
| `pwd` | Shows current path |
| `clear` | Clears the terminal |

### Content

| Command | Output |
|---------|--------|
| `cat help` | Full command reference |
| `cat profile` | Player stats — XP, Elo, level, badges, streak |
| `cat leaderboard` | Global top-20 rendered as a table in terminal |
| `cat daily` | Today's daily challenge brief |
| `cat missions/README` | Campaign list with status indicators |
| `cat missions/antariksha/README` | Campaign overview + chapter map (ASCII art India map) |
| `whoami` | Username + Elo band + rank |
| `history` | Recent command history within the app shell |

### Mission Commands (inside a chapter directory)

| Command | Output |
|---------|--------|
| `cat brief` | Displays the classified mission brief from story.md |
| `cat objectives` | Lists what needs to be done to complete the chapter |
| `start` | Launches the sandboxed terminal for this chapter |
| `submit <flag>` | Submits a flag — validates and awards XP |
| `hint` | Uses a hint (costs Elo) — reveals one clue |
| `status` | Shows chapter completion status |

### Easter Eggs / Special Commands

| Command | Output |
|---------|--------|
| `man shiva` | Lore entry on the SHIVA AI — in `man` page format |
| `man arjun` | Player's own dossier — updates as they level up |
| `grep flag *` | Returns: `"Nice try. Earn it. — SHIVA"` |
| `sudo su` | Returns: `"You don't have the clearance yet."` — unlocks at level 10 |
| `ping shiva` | Returns: `"Request timeout. SHIVA is not responding. Yet."` |
| `ssh shiva@satellite` | Returns the Chapter 10 teaser (only if ch1–9 complete) |
| `chmod 777 shiva` | Returns: `"Permission denied."` |

### Error Handling

Unknown commands return real bash-style errors:
```
guest@linuxquest:~$ fly
bash: fly: command not found
```

After 3 unknown commands in a row, a hint appears:
```
[CIRT-HINT] Type  cat help  to see available commands.
```

---

## Chapter Entry Flow

When a player runs `cd ch1`:

```
guest@linuxquest:~/missions/antariksha$ cd ch1

[CLASSIFIED TRANSMISSION INCOMING...]
████████████████████ 100%

arjun@sac-blr-01:~/missions/antariksha/ch1$ cat brief
```

The mission brief from `story.md` prints line by line with a
15ms typing delay — like a transmission being received.

Then:
```
arjun@sac-blr-01:~/missions/antariksha/ch1$ _
```

Two modes available from here:
1. **Read mode** — `cat brief`, `cat objectives`, `cat hint`
2. **Play mode** — `start` → opens the sandbox terminal as a split pane or full-screen takeover

---

## Terminal Layout

```
┌─────────────────────────────────────────────────────────────┐
│  [ISRO-CIRT]  ch1 — The Lab  |  Bangalore  |  ⚡ 1,247 XP  │  ← status bar (top, minimal)
├─────────────────────────────────────────────────────────────┤
│                                                             │
│                                                             │
│   (terminal output fills this space)                        │
│                                                             │
│                                                             │
│                                                             │
│                                                             │
│                                                             │
│  arjun@linuxquest:~$ _                                      │  ← prompt always at bottom
└─────────────────────────────────────────────────────────────┘
```

- **Status bar** (top): current location, chapter name, XP — single line, minimal
- **Body**: pure terminal output — no chrome, no borders, no cards
- **Prompt** (bottom): always pinned to the bottom

---

## Sandbox Mode (inside `start`)

When the player runs `start` inside a chapter:

```
arjun@sac-blr-01:~/missions/antariksha/ch1$ start

[SANDBOX INITIALIZING — sac-blr-01.isro.local]
Container ID: a3f9b2d1
Filesystem mounted: read-only (except /tmp, /home/player)
Network: isolated

You are now connected to the compromised server.
Type  exit  to return to mission shell.

arjun@sac-blr-01:/$ _
```

- This is a real Docker container via WebSocket
- `exit` or `Ctrl+D` returns to the mission shell (not the browser)
- Chapter sandbox has its own prompt (different hostname makes it clear)

---

## ASCII India Map (for `cat missions/antariksha/README`)

```
       Operation Antariksha — Mission Map

         ___
        /   \  ← Ch1 Bangalore    [✓ COMPLETE]
       | BLR |
        \___/
          |
       ___↓___
      /       \  ← Ch2 Chennai    [✓ COMPLETE]
     |   CHN   |
      \_______/
          |
       ___↓___
      /       \  ← Ch3 Mumbai     [► ACTIVE]
     |   MUM   |
      \_______/
          |
       ...etc
```

Completed cities show `✓`. Active city blinks. Locked cities show `?`.

---

## Visual Design

| Property | Value |
|----------|-------|
| Background | `#0d1117` (near-black) |
| Text | `#c9d1d9` (soft white) |
| Prompt | `#39ff14` (neon green) |
| Error | `#ff4444` (red) |
| Hint | `#ffaa00` (amber) |
| Success | `#39ff14` (green) |
| Locked/dim | `#444c56` (gray) |
| Font | `JetBrains Mono` or `Fira Code` — monospace only |
| Cursor | Blinking block `█` |

No other fonts. No colors outside this palette.

---

## Tech Stack for Frontend

| Layer | Choice | Why |
|-------|--------|-----|
| Framework | React + Vite (TypeScript) | Fast dev, easy WebSocket integration |
| Terminal emulator | `xterm.js` + `xterm-addon-fit` | Industry standard, handles ANSI, resizing |
| WebSocket | Native browser WebSocket | Connects to Go backend sandbox sessions |
| State | Zustand | Minimal, no redux overhead |
| Routing | None — single page, shell handles navigation | Stays true to shell-only UX |

---

## What NOT to Build

- ❌ No React Router pages
- ❌ No navbar, sidebar, or hamburger menu
- ❌ No modals or dialogs
- ❌ No buttons (except the initial "Connect" if needed for audio/focus)
- ❌ No cards, grids, or dashboard layouts
- ❌ No loading spinners — use terminal-style progress bars (`████████ 80%`)
- ❌ No toast notifications — print to terminal instead
- ❌ **No copy-paste** — block Ctrl+V, Ctrl+Shift+V, right-click, and the browser paste event

**Copy-paste implementation (xterm.js):**
```js
// Block Ctrl+V / Ctrl+Shift+V
terminal.attachCustomKeyEventHandler((e) => {
  if (e.ctrlKey && (e.key === 'v' || e.key === 'V')) {
    terminal.writeln('\r\n\x1b[33m[PASTE DISABLED] Type the command. No shortcuts here.\x1b[0m');
    return false;
  }
  return true;
});

// Block right-click context menu
terminal.element.addEventListener('contextmenu', (e) => e.preventDefault());

// Block browser-level paste
terminal.element.addEventListener('paste', (e) => {
  e.preventDefault();
  e.stopPropagation();
  terminal.writeln('\r\n\x1b[33m[PASTE DISABLED] Type the command. No shortcuts here.\x1b[0m');
});
```

---

## Login (Deferred)

Login is skipped for now. The shell starts as `guest@linuxquest`.
When implemented later, it will be done entirely in the terminal:

```
guest@linuxquest:~$ login

Username: arjun
Password: ████████

Authenticating...
Welcome back, Arjun. Elo: 1,742 | Rank: Sysadmin

arjun@linuxquest:~$ _
```

No forms. No redirect. Just shell prompts.
