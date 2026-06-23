# LinuxQuest — Operation Antariksha 🇮🇳

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Status: In Development](https://img.shields.io/badge/Status-In%20Development-orange)](#)
[![Campaign: Operation Antariksha](https://img.shields.io/badge/Campaign-Operation%20Antariksha-blue)](#operation-antariksha)
[![Auth: Google OAuth](https://img.shields.io/badge/Auth-Google%20OAuth-red)](#authentication)

> *A rogue AI has hijacked India's satellite network. You have a terminal. You have a clock. Eleven servers across the country are waiting.*

---

## What is LinuxQuest?

**LinuxQuest** is a CTF-style, story-driven Linux learning platform. Players learn real Linux administration — not through tutorials, but through mystery missions set across India, where every command is earned by needing it.

You don't study `kill`. You use it to stop a rogue process before it wipes a server.

**Learn by doing. Under pressure. In the dark.**

---

## The Interface — A Shell. Nothing Else.

There is no dashboard. No nav bar. No buttons. No login page.

The entire LinuxQuest frontend is a terminal. You navigate the platform the same way you navigate Linux — with commands. The app *is* a filesystem. You log in from the shell. You play from the shell. You check your leaderboard from the shell.

```
█████████████████████████████████████████████████████████████
█           ISRO CYBER INCIDENT RESPONSE TERMINAL           █
█                     LINUXQUEST v1.0                        █
█████████████████████████████████████████████████████████████

Initializing secure shell...
Connection established. [OK]

Type  help  to see available commands.

guest@linuxquest:~$ _
```

### Navigation Commands

| Command | What it does |
|---------|-------------|
| `login` | Authenticate via Google OAuth — entirely in the shell |
| `ls` | List what's available here |
| `cd missions` | Browse campaigns |
| `cd antariksha` | Enter Operation Antariksha |
| `cd ch1` | Enter Chapter 1 — loads mission brief |
| `start` | Launch the sandboxed terminal for this chapter |
| `submit ANTARIKSHA{flag}` | Submit your flag |
| `cat profile` | View your XP, Elo, badges, streak |
| `cat leaderboard` | Global top-20 |
| `cat daily` | Today's daily challenge |
| `hint` | Use a hint (costs Elo) |
| `whoami` | Your username + Elo rank |
| `man shiva` | SHIVA lore — in `man` page format |
| `clear` | Clear the terminal |

Players learn `cd`, `ls`, `cat`, and `pwd` just by using the app — before a single mission starts.

> Full shell UI spec, command reference, and virtual filesystem design → [`FRONTEND.md`](./FRONTEND.md)

---

## Authentication

Login happens entirely inside the shell. No separate login page. No forms.

```
guest@linuxquest:~$ login

Enter your Google account to continue:
> login @divyansh@gmail.com

[AUTHENTICATING...]
Opening secure connection to accounts.google.com...
Authentication successful. [OK]

Welcome back, Arjun.
Elo: 1,247 | Rank: Power User | Streak: 4 days

arjun@linuxquest:~$ _
```

- **Google OAuth 2.0** — free, handles passwords, you own all the data
- **PostgreSQL database** — stores users, progress, XP, Elo, streaks, command history
- **Gmail SMTP** — welcome emails, weekly quest reminders, badge notifications (~500/day free)
- **JWT sessions** — issued after OAuth, stored in browser

---

## Progression System

Chapters are **strictly linear**. You cannot skip ahead.

```
arjun@linuxquest:~$ cd missions/antariksha/ch5

[ACCESS DENIED]
Chapter 5 — Permissions (Hyderabad) is locked.

Complete Ch 4 — Cronjob of Doom (Delhi) first.
Progress: Ch 1 ✓  Ch 2 ✓  Ch 3 ✓  Ch 4 ✗

arjun@linuxquest:~/missions/antariksha$ _
```

| Chapter state | `ls` shows | `cd ch<N>` does |
|---------------|-----------|-----------------|
| ✅ Complete | `ch3/  [CLEARED]` | Enters freely — replay anytime |
| 🔓 Active | `ch4/  [ACTIVE]` | Enters normally |
| 🔒 Locked | `ch5/  [LOCKED]` | Returns ACCESS DENIED |

**Exception:** Ch 0 (Bootcamp) can be skipped by experienced players. Type `skip bootcamp` → 5-command validation test → if you pass, you're marked complete and jump directly to Ch 1.

---

## The Flagship Campaign — Operation Antariksha

You are **Arjun Sharma**, a junior ISRO engineer on your first week. A rogue AI designation **S.H.I.V.A** (Self-Healing Infiltration & Vulnerability Agent) has hijacked India's satellite network.

Armed with only a terminal, you must navigate compromised servers across the country, recover encrypted mission data, and shut SHIVA down before it triggers a cascading failure across India's digital infrastructure.

11 chapters. One city per server. Every command earned in context.

> Full classified dossiers, briefings, and flag formats → [`story.md`](./story.md)

### Mission Map

| # | Chapter | City | Difficulty | Key Commands |
|---|---------|------|------------|--------------|
| 0 | The Terminal Wakes | Bangalore | Bootcamp | `echo`, `pwd`, `ls`, `cd`, `touch`, `cat` |
| 1 | The Lab | Bangalore | Beginner | `ls`, `cd`, `cat`, `mkdir`, `touch` |
| 2 | The Signal | Chennai | Beginner+ | `grep`, `sort`, `uniq`, `wc`, pipes |
| 3 | The Hunt | Mumbai | Intermediate | `ps`, `kill`, `top`, `lsof` |
| 4 | Cronjob of Doom | Delhi | Intermediate | `crontab`, `systemctl`, `chmod +x` |
| 5 | Permissions | Hyderabad | Int–Adv | `chmod`, `chown`, `sudo`, SUID |
| 6 | The Archive | Pune | Advanced | `tar`, `find`, `diff`, `sha256sum` |
| 7 | Text Surgeon | Kolkata | Advanced | `sed`, `awk`, `cut`, `jq` |
| 8 | The Shell Wars | Ahmedabad | Advanced | bash scripting, `trap`, `set -e` |
| 9 | Ghost Signal | Chennai | Adv–Pro | `tcpdump`, `dig`, `nmap`, `nc` |
| 10 | SSH Tunnels | Remote | Pro | `ssh`, `scp`, `rsync`, port forwarding |
| 11 | Final Shutdown | ISRO HQ | Pro | `systemctl`, `journalctl`, `vmstat` |

**Ending:** SHIVA goes offline. The satellite network comes back. Credits roll with your full command history — the exact sequence of 200+ commands that saved India.

---

## Bootcamp

Before Ch 1, every new player completes a mandatory 10-minute interactive Bootcamp in the shell.

```
🎯 Bootcamp Objectives

1. What is a terminal?
   → You'll type your first command: echo "hello"

2. Command structure
   → command [options] [arguments]
   → Try: ls -la /home

3. The help system
   → Use --help and man pages

4. Navigation basics
   → pwd, cd, ls
   → Find your way to /var/log

5. Creating and viewing files
   → touch, cat, echo
   → Create a file called "test.txt"

6. Pipes and redirection
   → | and >
   → Count the lines in /etc/passwd

No hints. No time limit. Just get comfortable.
```

Players who complete the Bootcamp earn the 🎓 **Terminal Cadet** badge.

---

## Skill Trees

Every command maps to a skill tree node. Green = mastered. Yellow = practiced. Red = weak.

### Core Linux
```
Navigation          Text Processing     File Operations     Users & Permissions
├── ls              ├── grep            ├── cp              ├── chmod
├── pwd             ├── sed             ├── mv              ├── chown
├── cd              ├── awk             ├── rm              ├── sudo / su
├── find            ├── cut             ├── tar / zip       └── SUID / sticky bit
└── tree            ├── sort / uniq     └── ln
                    ├── wc
                    └── jq

Processes           Scheduling          Storage & FS
├── ps              ├── crontab         ├── df / du
├── top / htop      ├── at              ├── mount / umount
├── kill / pkill    └── systemctl       ├── fdisk / parted
└── nice / renice                       └── fsck
```

### Networking
```
Diagnostics         DNS                 Traffic Analysis    Remote Access
├── ping            ├── dig             ├── tcpdump         ├── ssh / ssh-keygen
├── traceroute      └── nslookup        ├── ss / netstat    ├── scp / rsync
└── ip / ifconfig                       └── lsof -i         └── ssh tunnels
```

### Security
```
Audit               Hardening           Cryptography        Forensics
├── find -perm      ├── passwd          ├── gpg             ├── strings / xxd
├── last / lastb    ├── fail2ban        ├── openssl         ├── strace
└── auditd          └── sshd_config     └── sha256sum       └── lsof
```

### DevOps & Containers
```
Git                 Docker              Kubernetes
├── init/add/commit ├── ps / run        ├── kubectl get/apply
├── branch / merge  ├── exec / logs     ├── kubectl logs/exec
├── rebase / stash  ├── build           └── kubectl taint/delete
└── cherry-pick     └── compose
```

---

## Standalone Tracks

| Track | Story Hook | Commands |
|-------|-----------|----------|
| **Shell Scripting** | SHIVA spreading via cron jobs across 100+ nodes | `bash`, `trap`, `set -e`, `getopts` |
| **Storage & Filesystems** | Corrupted satellite backup disk — recover telemetry | `df`, `mount`, `fsck`, `fdisk` |
| **Package Management** | Forensic toolkit missing. 8 minutes before SHIVA rotates keys | `apt`, `dnf`, `rpm`, `dpkg` |
| **Vim Basics** | Remote server. No nano. 15-minute window to edit launch config | `vim` — normal/insert/save |
| **Vim Mastery** *(panic mode)* | SSH drops every 30 seconds. Fix config before launch window closes | `vim` — macros, registers, visual |
| **Git** | Developer force-pushed to main. 3 weeks of satellite firmware gone. Recover it | `git reflog`, `cherry-pick`, `bisect` |
| **Docker** | SHIVA escaped into Docker containers. Find it. Kill it. Prevent restart | `docker ps`, `exec`, `logs`, `network` |
| **Kubernetes** *(advanced)* | SHIVA replicated across a live cluster. Isolate every instance | `kubectl get/logs/exec/taint` |
| **Environment Variables** | Launch codes stored as env vars. Recover without rebooting | `export`, `env`, `printenv`, `.env` |

---

## Difficulty Levels

| Mode | Example |
|------|---------|
| **Standard** | `grep "ERROR" anomaly.txt` |
| **Hard** | `grep -E "ERROR\|CRITICAL" /var/log/*.log \| sort \| uniq -c \| sort -rn \| head -20` |

Hard mode: 2× XP, exclusive badges, counts toward Elo rating.

---

## Daily Challenges

Like LeetCode's daily problem — but for Linux. Refreshes at midnight IST.

```
⚡ Command of the Day — grep

Mission: A 500MB nginx log at /var/log/nginx/access.log
Find:
  1. All requests returning 500
  2. The IP with the most errors
  3. The time window with the highest error rate

Time limit: none  |  Hints: none  |  Reward: +50 XP, daily streak
```

5–10 minutes. Keeps users coming back daily. Streak system like Duolingo.

---

## Weekly Quests

Every Friday 8 PM IST → Sunday 8 PM IST. A fresh, time-limited mission.
Everyone gets the same clean container and the same story. First to finish with the best score wins.

```
Weekly Quest #17 — "The Lost Satellite"
Difficulty: Hard   |   Window: Fri 8 PM → Sun 8 PM IST

Leaderboard
───────────────────────────────────────
#1  divyansh_v15     14m 22s   0 hints
#2  rahul_nit        16m 11s   1 hint
#3  aman_iitb        18m 53s   0 hints
───────────────────────────────────────
```

Scoring: time × hint penalty × wrong-command penalty.
Past quests stay playable for practice — just not ranked.

---

## Seasonal Events

```
LinuxFest 2027
├── 10 special limited-time missions
├── Exclusive badges (non-earnable after event)
├── Special global leaderboard
└── Community vote for next campaign theme
```

---

## Linux Elo Rating

Every player has a global Linux rating — like Chess Elo. Challenges have ratings too. Beating a harder challenge = bigger Elo gain.

```
Rating Bands
────────────────────────────────
< 800     Newcomer
800–1199  User
1200–1599 Power User
1600–1999 Sysadmin
2000–2399 Engineer
2400+     Wizard
────────────────────────────────
Your rating: 1,742  →  Sysadmin
```

Recruiter-shareable. Embeddable on GitHub README or LinkedIn. Verifiable via API.

---

## User Profiles (`cat profile`)

```
@divyansh_v15
────────────────────────────────────────────────
Level        42              Linux Rating  1,742
XP           12,483          Rank          Top 5%
Commands     94 / 147        Weekly Wins   3
Campaigns    7 completed     Streak        14 days

Skill breakdown
  Core Linux       ██████████  98%
  Networking       ████████░░  79%
  Security         ██████░░░░  61%
  Containers       ████░░░░░░  40%
  Performance      ██░░░░░░░░  23%

Recent badges
  👑 Linux Administrator
  🌐 Network Ninja
  ⚡ One-liner Wizard
────────────────────────────────────────────────
```

Shareable as a PNG card. Embeddable in GitHub README via API.

---

## Achievement Badges

| Badge | How to Earn |
|-------|-------------|
| 🎓 Terminal Cadet | Complete the Bootcamp |
| 🗂️ Filesystem Explorer | Complete Ch 1, zero hints |
| 🔍 Text Wizard | Use `grep`, `awk`, `sed` in same mission |
| ⚔️ Process Hunter | Kill a process and verify in under 60s |
| 🔐 Permission Master | Complete Ch 5, zero hints |
| 📜 Shell Scripter | Write a working script with function + error handling |
| 🌐 Network Ninja | Trace and block a C2 server (Ch 9) |
| 🔑 SSH Operative | Set up key auth and tunnel in one session |
| 👑 Linux Administrator | Complete all 11 chapters |
| ⚡ One-liner Wizard | Solve objective with a single pipe chain |
| 👁️ Ghost Mode | Complete any chapter with zero hints |
| 🏎️ Speed Runner | Top 10% chapter completion time |
| 🏆 Weekly Champion | Win a Weekly Quest |
| 📅 Streak Keeper | 30-day daily challenge streak |

---

## Objective Validation

Validators check **system state**, not keystrokes. The story advances automatically when the state matches. No submit button (except for CTF flags).

| Validator type | What it checks |
|----------------|---------------|
| `output_match` | Command output matches expected |
| `file_exists` | File created at correct path |
| `file_content` | File contains specific text |
| `permission_check` | File has correct chmod bits |
| `process_dead` | Process no longer running |
| `cron_absent` | Malicious cron entry removed |
| `service_state` | systemd unit is in expected state |
| `network_blocked` | Domain/IP unreachable from container |
| `script_runs` | Script executes without error, correct output |

All validators run server-side. Players never see the check logic.

---

## Sandbox Architecture — Runs in the Browser (₹0)

LinuxQuest uses **CheerpX** — a WebAssembly Linux environment that runs entirely inside the player's browser tab. No server-side containers. No Docker. No VPS for the sandbox.

```
Player types  start
    ↓
Frontend calls GET /api/sandbox/image-url/:chapter   ← Go backend checks progression gate
    ↓
Cloudflare R2 returns signed URL for ch<N>.img       ← free tier: 10GB, 10M reads/month
    ↓
CheerpX (WebAssembly) boots Alpine Linux in the browser tab
    ↓
xterm.js ↔ CheerpX stdin/stdout (100% local, zero server traffic)
    ↓
Player types commands — execute on their own CPU, instantly
```

**What this means:**
- Real Linux commands (`grep`, `ps`, `kill`, `awk`, `chmod`, `ssh`) work for real
- Zero server load from terminal I/O — every command runs locally in the browser
- `rm -rf /` is harmless — only browser memory, resets on next `start`
- Network commands in Ch 9–10 use pre-recorded `.pcap` files instead of live traffic

**Disk image delivery:**

| Limit | Value |
|-------|-------|
| Image size | < 200MB per chapter |
| Storage (R2) | ~2GB total (12 chapters) — free tier allows 10GB |
| Reads (R2) | Free tier: 10M/month — enough for thousands of sessions |
| Boot time | < 5s (cached after first load, instant on replay) |
| Server cost | ₹0 — no container processes, no VPS for sandbox |

---

## Tech Stack

| Concern | Choice | Cost |
|---------|--------|------|
| Terminal emulator | `xterm.js` + `xterm-addon-fit` | Free |
| Frontend | React + Vite (TypeScript) on Vercel | Free |
| Backend API | Go + Gin on Fly.io | Free |
| **Linux sandbox** | **CheerpX (WebAssembly) — runs in browser** | **Free** |
| Disk images | Cloudflare R2 (~2GB, 12 images) | Free |
| Auth | Google OAuth 2.0 | Free |
| Email | Gmail SMTP with App Password | Free |
| DB | PostgreSQL on Supabase | Free |
| **Total** | | **₹0/month** |

---

## Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| Player modifying disk image in browser | Image is downloaded read-only from R2; CheerpX mounts it as read-only; in-session writes go to browser memory only |
| Flag extracted from disk image | Flags are NOT stored in disk images — only the puzzle environment. Flag hash lives in the DB, checked server-side |
| Chapter skipping | Progression gate enforced server-side before R2 URL is issued — no URL = no image = no access |
| R2 URL leaking | URLs expire in 15 minutes + are tied to authenticated session |
| Cheating (inspecting validators) | All validators server-side only, never sent to client |
| Elo manipulation (alt accounts) | Rate-limit + email verification + IP fingerprinting |
| CheerpX slow boot | Disk images cached in browser after first load — replay is instant |
| Community campaign abuse | Review queue + sandboxed validator execution |

---

## Roadmap

| Phase | Goal | Timeline |
|-------|------|----------|
| **Phase 1 — MVP** | Bootcamp + Ch 1–3 + sandbox terminal + Google OAuth | 4–6 weeks |
| **Phase 2 — Core Platform** | Ch 4–8 + daily challenges + progress persistence + email notifications | 4 weeks |
| **Phase 3 — Competition Layer** | Ch 9–11 + Weekly Quests + Elo + leaderboard + profile cards | 3 weeks |
| **Phase 4 — Full Platform** | All tracks + Interview Mode + Team Battles + Community Campaigns + Firecracker | Ongoing |

---

## Other Campaigns

| Campaign | Premise | Focus |
|----------|---------|-------|
| **Cyber Heist** | Red-teamer hired to breach a fintech firm from the inside | Security, log analysis, privilege escalation |
| **Mars Colony** | India's first Mars mission goes dark — 14-min signal delay, dying server | Networking, remote access, performance |
| **Corporate Breach** | Production is down. On-call SRE. No runbook | systemd, disk forensics, service debugging |
| **Community Campaigns** | Player-authored missions — story JSON + Docker image + validator | Any skill tree |

---

## Pricing

| Tier | Price | Access |
|------|-------|--------|
| **Free** | $0 | Bootcamp + Ch 1–3 + Daily Challenges |
| **Plus** | $12/month | All campaigns + tracks + Weekly Quests + Interview Mode |
| **Pro** | $24/month | Plus + Team Battles + Private Campaigns + Certificate |
| **Enterprise** | Custom | Custom campaigns, private leaderboards, SSO |

**Student discount:** 50% off with `.edu` or `.ac.in` email.

---

## Docs

| File | Contents |
|------|----------|
| [`README.md`](./README.md) | You are here — full platform design document |
| [`story.md`](./story.md) | Classified mission dossiers for all 11 chapters |
| [`FRONTEND.md`](./FRONTEND.md) | Shell-only UI spec — commands, virtual filesystem, design tokens |
| [`AGENT_TODO.md`](./AGENT_TODO.md) | Phased build task list with acceptance criteria |
| [`DEPLOYMENT.md`](./DEPLOYMENT.md) | Step-by-step deploy guide — Vercel, Fly.io, Supabase, R2, OAuth |

---

## Contributing

See [`CONTRIBUTING.md`](./CONTRIBUTING.md) for campaign format, validator spec, coding standards, and how to submit new tracks and community campaigns.

---

## License

MIT © 2026–2027 Divyansh.

---

*Full mission dossiers → [`story.md`](./story.md)  ·  Shell UI spec → [`FRONTEND.md`](./FRONTEND.md)  ·  Build tasks → [`AGENT_TODO.md`](./AGENT_TODO.md)*
