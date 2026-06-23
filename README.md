```markdown

\# LinuxQuest



> The Duolingo for Linux. A platform where players learn real Linux administration through story-driven missions, weekly competitions, and skill-based progression — from `ls` to `kubectl`.



\---



\## Why LinuxQuest?



Most Linux tutorials teach commands in isolation. Students read about `grep`, run one example, move on, and forget it in a week.



LinuxQuest teaches Linux the way you actually learn it — by needing it. Every command is introduced at the exact moment it's required to progress the story. You don't study `kill`; you use it to stop a rogue process before it wipes a server.



\*\*Learn by doing, not memorizing.\*\*



\---



\## What You'll Learn



By completing LinuxQuest, players will be able to:



\- Navigate and manipulate Linux filesystems confidently

\- Search and process text with `grep`, `awk`, `sed`, `jq`

\- Manage processes, services, cron jobs, and systemd units

\- Understand file permissions, ownership, and privilege escalation

\- Write defensive shell scripts with error handling

\- Trace, capture, and block network traffic

\- Use SSH, port forwarding, and remote file transfer

\- Manage packages, storage, and filesystems

\- Work with containers (Docker, Kubernetes)

\- Administer, monitor, and tune a production Linux system

\- Approach real-world DevOps and SRE tasks with confidence



\---



\## Platform Overview



LinuxQuest is not a single game — it is a Linux learning platform. Content is organized into \*\*Skill Trees\*\*, playable as \*\*Campaigns\*\* (story-driven) or \*\*Tracks\*\* (focused practice). Competitive players engage via \*\*Weekly Quests\*\* and a global \*\*Linux Elo Rating\*\*.



```

LinuxQuest Platform

│

├── 🌳 Skill Trees          Core Linux · Networking · Security · DevOps

│                           Containers · Cloud · Databases · Incident Response

│                           Performance Tuning · Git

│

├── 🎮 Campaigns            Operation Antariksha  (flagship, ISRO theme)

│                           Cyber Heist · Mars Colony · Corporate Breach

│                           + Community-created campaigns

│

├── ⚡ Daily Challenges      Command of the Day — 5 min mission

│

├── 🏆 Weekly Quests        Global leaderboard · Time-limited · Fresh containers

│

├── 🎓 Interview Mode        Real server tasks, no story, no hints

│

└── 👥 Team Battles          College vs college · Club competitions · Hackathons

```



\---



\## Skill Trees



Every command in LinuxQuest maps to a skill tree node. The tree is the player's roadmap — green = mastered, yellow = practiced, red = weak.



\### Core Linux

```

Navigation          Text Processing     File Operations     Users \& Permissions

├── ls              ├── grep            ├── cp              ├── chmod

├── pwd             ├── sed             ├── mv              ├── chown

├── cd              ├── awk             ├── rm              ├── sudo / su

├── find            ├── cut             ├── tar / zip       ├── visudo

└── tree            ├── sort / uniq     ├── diff            └── SUID / sticky bit

&#x20;                   ├── wc              └── ln

&#x20;                   └── jq



Processes           Scheduling          Environment         Storage \& FS

├── ps              ├── crontab         ├── export          ├── df / du

├── top / htop      ├── cron.d          ├── env             ├── lsblk

├── kill / pkill    ├── at              ├── printenv        ├── mount / umount

├── nice / renice   └── systemctl       ├── unset           ├── fdisk / parted

└── jobs / bg / fg    list-timers      └── .bashrc / .env  ├── mkfs

&#x20;                                                           ├── fsck

&#x20;                                                           └── blkid

```



\### Networking

```

Diagnostics         DNS                 Traffic Analysis    Firewall

├── ping            ├── dig             ├── tcpdump         ├── iptables

├── traceroute      ├── nslookup        ├── wireshark       ├── ufw

├── mtr             └── host            ├── ss / netstat    └── nftables

└── ip / ifconfig                       └── lsof -i



Remote Access       Transfer            Scanning

├── ssh             ├── scp             ├── nmap

├── ssh-keygen      ├── rsync           └── nc (netcat)

├── ssh tunnels     ├── curl

└── tmux / screen   └── wget

```



\### Security

```

Audit               Hardening           Cryptography        Forensics

├── find -perm      ├── passwd          ├── gpg             ├── strings

├── last / lastb    ├── usermod         ├── openssl         ├── xxd

├── auditd          ├── fail2ban        ├── md5sum          ├── strace

└── lynis           └── sshd\_config     └── sha256sum       └── lsof

```



\### DevOps \& Containers

```

Git                 Docker              Kubernetes          CI/CD

├── init/add/commit ├── ps / run        ├── kubectl get     ├── Dockerfile

├── branch / merge  ├── exec / logs     ├── kubectl apply   ├── docker-compose

├── rebase / stash  ├── build           ├── kubectl logs    └── systemd units

└── cherry-pick     └── compose         └── kubectl exec

```



\### Performance Tuning

```

CPU \& Memory        Disk I/O            Network             Profiling

├── vmstat          ├── iostat          ├── iftop           ├── perf

├── sar             ├── iotop           ├── nethogs         ├── strace

├── free -h         ├── hdparm          └── ss -s           └── ltrace

└── dmesg           └── dd

```



\---



\## The Bootcamp



Before diving into Chapter 1, every new player completes a mandatory 10‑minute interactive Bootcamp. This ensures no one gets stuck on terminal basics.



```

🎯 Bootcamp Objectives



1\. What is a terminal?

&#x20;  → You'll type your first command: echo "hello"



2\. Command structure

&#x20;  → command \[options] \[arguments]

&#x20;  → Try: ls -la /home



3\. The help system

&#x20;  → Use --help and man pages

&#x20;  → Find the option for "human-readable" sizes in ls



4\. Navigation basics

&#x20;  → pwd, cd, ls

&#x20;  → Find your way to /var/log



5\. Creating and viewing files

&#x20;  → touch, cat, echo

&#x20;  → Create a file called "test.txt" with content "Hello Bootcamp"



6\. Pipes and redirection

&#x20;  → | and >

&#x20;  → Count the lines in /etc/passwd



No hints. No time limit. Just get comfortable.

```



Players who complete the Bootcamp earn the 🎓 \*\*Terminal Cadet\*\* badge. Skipping the Bootcamp isn't allowed — it's the foundation for everything that follows.



\---



\## Campaigns



\### Flagship — Operation Antariksha 🇮🇳



You are \*\*Arjun\*\*, a young ISRO engineer. On your first week, a rogue AI named SHIVA hijacks India's satellite network. Armed with only a terminal, you must navigate compromised servers across the country, recover encrypted mission data, and shut SHIVA down before it triggers a cascading failure across India's digital infrastructure.



11 chapters. One city per server. Every command earned in context.



| # | Chapter | City | Difficulty | Key Commands |

|---|---|---|---|---|

| 1 | The Lab | Bangalore | Beginner | `ls`, `cd`, `cat`, `mkdir`, `touch` |

| 2 | The Signal | Chennai | Beginner+ | `grep`, `sort`, `uniq`, `wc`, pipes |

| 3 | The Hunt | Mumbai | Intermediate | `ps`, `kill`, `top`, `lsof` |

| 4 | Cronjob of Doom | Delhi | Intermediate | `crontab`, `systemctl`, `chmod +x` |

| 5 | Permissions | Hyderabad | Int–Adv | `chmod`, `chown`, `sudo`, SUID |

| 6 | The Archive | Pune | Advanced | `tar`, `find`, `diff`, `sha256sum` |

| 7 | Text Surgeon | Kolkata | Advanced | `sed`, `awk`, `cut`, `jq` |

| 8 | The Shell Wars | Ahmedabad | Advanced | bash scripting, `trap`, `set -e` |

| 9 | Ghost Signal | Chennai-2 | Adv–Pro | `tcpdump`, `dig`, `nmap`, `nc` |

| 10 | SSH Tunnels | Remote | Pro | `ssh`, `scp`, `rsync`, port forwarding |

| 11 | Final Shutdown | ISRO HQ | Pro | `systemctl`, `journalctl`, `vmstat` |



\*\*Ending:\*\* SHIVA goes offline. The satellite network comes back. Credits roll with your full command history — the exact sequence of 200+ commands that saved India.



\---



\### Other Campaigns



\*\*Cyber Heist\*\* — Corporate espionage. You're a red-teamer hired to test a fintech firm's defenses. Focuses on security, log analysis, and privilege escalation.



\*\*Mars Colony\*\* — A comms blackout on India's first Mars mission. Focuses on networking, remote access, and performance debugging under degraded conditions.



\*\*Corporate Breach\*\* — Incident response. You're the SRE on-call when production goes down. Focuses on systemd, logs, disk, and service debugging.



\*\*Community Campaigns\*\* — Anyone can author and publish a campaign: story JSON + Docker image + validator config. Rated, upvoted, and searchable.



> \*\*Community Moderation:\*\* Campaigns go through a review queue before publication. Malicious validators run in isolated sandboxes. Community flagging and peer review keep quality high.



\---



\## Standalone Tracks



These concepts don't fit cleanly into Operation Antariksha but are essential for a complete Linux education. They ship as dedicated tracks with their own story framing.



\### Shell Scripting Track \*(NEW)\*

\*\*Story:\*\* SHIVA is spreading through cron jobs. You need to write defensive scripts to detect, report, and neutralize the threat across 100+ nodes — without manual intervention.



| Chapter | Topic |

|---|---|

| 1 | Variables, loops, conditionals |

| 2 | Functions and error handling (`set -e`, `trap`) |

| 3 | Parsing command-line arguments (`getopts`) |

| 4 | Working with JSON and CSV (`jq`, `awk`) |

| 5 | Cron integration and rotating logs |

| 6 | Systemd service management from scripts |



This is the bridge between "knowing commands" and "automating systems" — a critical gap in most Linux curricula.



\### Storage \& Filesystems Track

\*\*Story:\*\* The AI corrupted the satellite's backup disk. Mount it, repair the filesystem, recover the telemetry data before the ground window closes.



Commands: `df`, `du`, `lsblk`, `mount`, `umount`, `blkid`, `fdisk`, `parted`, `mkfs`, `fsck`



\### Package Management Track

\*\*Story:\*\* Your forensic toolkit isn't installed on the compromised node. You have 8 minutes before the AI rotates its keys.



Ubuntu path: `apt`, `apt-cache`, `apt-get`, `dpkg`

Fedora path: `dnf`, `rpm`



\### Vim Track



The Vim track is split into two tiers to avoid the "panic mode" overwhelm:



\#### Vim Basics Track

\*\*Story:\*\* The remote server has no nano, no VSCode, no internet. You have 15 minutes to edit the launch configuration before the window closes. The SSH connection is stable — but the clock is ticking.



Commands: `vim` — normal mode, insert mode, save/quit (`:wq`), search (`/pattern`), replace (`:%s/old/new/g`)



\#### Vim Mastery Track \*(panic mode)\*

\*\*Story:\*\* The SSH connection drops every 30 seconds. Fix the config before the satellite launch window closes. No reconnects. No retries.



Commands: `vim` — macros (`q`), registers (`"`), visual mode (`v`), multi-file editing (`:n`, `:prev`), advanced search/replace



> The Basics track is required. The Mastery track is for those who want the badge.



\### Git Track

\*\*Story:\*\* A developer accidentally force-pushed to main and destroyed 3 weeks of satellite firmware code. Recover it.



Commands: `git init`, `add`, `commit`, `branch`, `merge`, `rebase`, `stash`, `reflog`, `cherry-pick`, `bisect`



\### Docker Track

\*\*Story:\*\* SHIVA escaped the main server and is hiding inside Docker containers. Find it. Kill it. Prevent it from restarting.



Commands: `docker ps`, `docker run`, `docker logs`, `docker exec`, `docker inspect`, `docker compose`, `docker network`



\### Kubernetes Track \*(advanced)\*

\*\*Story:\*\* SHIVA replicated itself across a Kubernetes cluster. It keeps spawning new pods. You need to find, isolate, and terminate every instance — while the cluster is live.



Commands: `kubectl get`, `kubectl describe`, `kubectl logs`, `kubectl exec`, `kubectl delete`, `kubectl apply`, `kubectl taint`



\### Environment Variables Track

\*\*Story:\*\* The mission launch codes are stored as environment variables on a locked system. Recover them without rebooting.



Commands: `echo $VAR`, `env`, `export`, `unset`, `printenv`, `.env` files, `envsubst`



\---



\## Difficulty Levels



Every objective ships in two modes. Players choose before starting.



| Mode | Example |

|---|---|

| Standard | `grep "ERROR" anomaly.txt` |

| Hard | `grep -E "ERROR\\|CRITICAL" /var/log/\*.log \\| sort \\| uniq -c \\| sort -rn \\| head -20` |



Hard mode: 2× XP, unlocks exclusive badges, counts toward Elo rating.



\---



\## Daily Challenges



Like LeetCode's daily problem — but for Linux.



```

⚡ Command of the Day — grep



Mission: A 500MB nginx log dropped into /var/log/nginx/access.log

Find:

&#x20; 1. All requests returning 500

&#x20; 2. The IP address with the most errors

&#x20; 3. The time window with the highest error rate



Time limit: none

Hints: none

Reward: +50 XP, daily streak

```



5–10 minutes. Keeps users coming back daily. Streak system identical to Duolingo.



\---



\## Weekly Quests



Every Friday 8 PM IST → Sunday 8 PM IST. A fresh, time-limited mission. Everyone gets the same clean container and the same story. First to finish with the best score wins.



```

Weekly Quest #17 — "The Lost Satellite"

Difficulty: Hard

Window: Fri 8 PM → Sun 8 PM IST



Leaderboard

───────────────────────────────────────

\#1  divyansh\_v15     14m 22s   0 hints

\#2  rahul\_nit        16m 11s   1 hint

\#3  aman\_iitb        18m 53s   0 hints

\#4  priya\_dev        21m 04s   2 hints

───────────────────────────────────────

```



Scoring: time × hint penalty × wrong-command penalty.



Past quests stay playable for practice — just not ranked.



\---



\## Seasonal Events



```

LinuxFest 2027

├── 10 special limited-time missions

├── Exclusive badges (non-earnable after event)

├── Special global leaderboard

└── Community vote for next campaign theme

```



\---



\## Linux Elo Rating



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



\---



\## User Profiles



```

@divyansh\_v15

────────────────────────────────────────────────

Level        42              Linux Rating  1,742

XP           12,483          Rank          Top 5%

Commands     94 / 147        Weekly Wins   3

Campaigns    7 completed     Streak        14 days



Skill breakdown

&#x20; Core Linux       ██████████  98%

&#x20; Networking       ████████░░  79%

&#x20; Security         ██████░░░░  61%

&#x20; Containers       ████░░░░░░  40%

&#x20; Performance      ██░░░░░░░░  23%



Recent badges

&#x20; 👑 Linux Administrator

&#x20; 🌐 Network Ninja

&#x20; ⚡ One-liner Wizard

────────────────────────────────────────────────

```



Shareable as a card (PNG export). API endpoint for embedding in GitHub README.



\---



\## Command Graph



Every player's personal map of Linux mastery. Visible on their profile.



```

Navigation          Text Processing     Processes

├── ls      ✅      ├── grep    ✅      ├── ps      ✅

├── pwd     ✅      ├── sed     ✅      ├── kill    ✅

├── cd      ✅      ├── awk     🟡      ├── top     ✅

├── find    🟡      ├── cut     ✅      ├── htop    🟡

└── tree    🔴      ├── sort    ✅      └── nice    🔴

&#x20;                   └── jq      🔴

```



✅ Mastered (used correctly 5+ times)

🟡 Practiced (used but not consistently)

🔴 Weak or untouched



The graph is the player's learning roadmap. It shows exactly what's left — not "Chapter 6 locked" but "you haven't touched `awk` yet."



\---



\## Achievement System



Badges tied to real Linux skill categories:



| Badge | How to Earn |

|---|---|

| 🎓 Terminal Cadet | Complete the Bootcamp |

| 🗂️ Filesystem Explorer | Complete all Ch1 objectives, zero hints |

| 🔍 Text Wizard | Use `grep`, `awk`, and `sed` in the same mission |

| ⚔️ Process Hunter | Kill a process and verify it's gone in under 60s |

| 🔐 Permission Master | Complete Ch5, zero hint usage |

| 📜 Shell Scripter | Write a working script with a function and error handling |

| 🌐 Network Ninja | Trace and block a C2 server (Ch9) |

| 🔑 SSH Operative | Set up key auth and tunnel in one session |

| 👑 Linux Administrator | Complete all 11 chapters |

| ⚡ One-liner Wizard | Solve an objective with a single pipe chain |

| 👁️ Ghost Mode | Complete any chapter with zero hints |

| 🏎️ Speed Runner | Top 10% chapter completion time |

| 📦 Package Master | Install, verify, and pin a package without breaking deps |

| 🐳 Container Hunter | Find and kill a process hiding inside Docker |

| ☸️ Cluster Warden | Complete Kubernetes track |

| 📅 Streak Keeper | 30-day daily challenge streak |

| 🏆 Weekly Champion | Win a Weekly Quest |

| 📜 Bash Automator | Complete the Shell Scripting Track |

| ✍️ Vim Sage | Complete Vim Mastery Track (panic mode) |



\---



\## Command Line Hero



After completing Operation Antariksha, players unlock \*\*Command Line Hero\*\* — a cinematic replay of their entire journey.



```

────────────────────────────────────────────────

&#x20;        OPERATION ANTARIKSHA

&#x20;        ───────────────────



&#x20;   Your command history:

&#x20;   \[ls -la] → \[cd /var/log]

&#x20;   → \[grep -r "SHIVA" .]

&#x20;   → \[ps aux | grep shiva]

&#x20;   → \[kill -9 847]

&#x20;   → \[systemctl stop shiva.service]

&#x20;   → \[crontab -r]

&#x20;   → \[ssh -L 4444:localhost:4444 arjun@delhi-node]

&#x20;   → \[sed -i 's/password=plaintext/password=REDACTED/g' config.conf]

&#x20;   → \[kubectl delete pods -l app=shiva]

&#x20;   → \[systemctl start network]



&#x20;   147 commands.

&#x20;   11 chapters.

&#x20;   1 satellite saved.



&#x20;   ══════════════════════════════════════════════

&#x20;   YOU ARE NOW A LINUX ADMINISTRATOR.

&#x20;   ══════════════════════════════════════════════



&#x20;   Share this video? \[Y/n]

────────────────────────────────────────────────

```



The replay is generated from the player's saved command history and rendered as a scrolling animation. Shareable as a video or embed.



\---



\## Interview Mode



No story. No hints. No objective list. Real server tasks — exactly what companies ask in SRE and DevOps interviews.



```

Linux Interview Track — Problem #14



Your production web server is throwing 502s.

The server has been up for 3 days.



Diagnose and fix the issue.



You have access to:

&#x20; - Full root shell

&#x20; - systemd, nginx, postgres running (or not)

&#x20; - /var/log/\*



No hints. No time limit. Scored on:

&#x20; - Commands used (efficiency)

&#x20; - Time to resolution

&#x20; - Whether you actually fixed it

```



Scenarios include: disk full, process OOM-killed, cron job gone rogue, misconfigured nginx, zombie processes, and port conflicts.



\---



\## Team Battles



College vs college. Club vs club. Company onboarding cohorts.



```

NIT Hamirpur  vs  NIT Jalandhar

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Players      12                  11

Avg Time     18m 44s             22m 11s

Hints Used   14                  27

Objectives   94 / 110            81 / 110

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Winner: NIT Hamirpur 🏆

```



Use cases: hackathon side events, college Linux clubs, DevOps bootcamp cohorts, corporate onboarding.



\---



\## Command Handbook (In-Game)



Always accessible from a side drawer. Unlocks progressively as the player encounters each command.



```

📖 Command Handbook



── Unlocked (23) ──────────────────────────────

✓ grep      Search text patterns in files

✓ chmod     Change file permissions

✓ ps        Show running processes

✓ kill      Terminate a process by PID



── Examples ────────────────────────────────────

grep "ERROR" logs.txt

grep -r "shiva" /var/log/

grep -i "critical" \*.log | wc -l



── Locked (51) ─────────────────────────────────

&#x20; tcpdump   (unlocked in Chapter 9)

&#x20; kubectl   (unlocked in Kubernetes track)

```



\---



\## Analytics Dashboard



```

Your Linux Stats

────────────────────────────────────────────────

Commands learned       94 / 147

Unique commands used   81

Overall success rate   84%

Avg hints / mission     1.2

Total time played      14h 22m

Current streak         14 days



Skill breakdown

&#x20; Text Processing   ████████████  92%

&#x20; File Operations   ██████████░░  84%

&#x20; Processes         █████████░░░  76%

&#x20; Networking        ████░░░░░░░░  38%    ← weakest

&#x20; Containers        ███░░░░░░░░░  28%



Most used command   grep (247 times)

Best solve          Ch7 Obj3 — awk + sort + uniq -c in one line

────────────────────────────────────────────────

```



\---



\## Real World Mode



After each chapter, an unscripted challenge unlocks. No story. No hints. No objective list. Just a problem and a terminal:



```

Real World Challenge — Chapter 2



You have a 2.3 GB nginx access log at /var/log/nginx/access.log



Find:

&#x20; 1. The IP address with the most requests

&#x20; 2. Total number of 404 responses

&#x20; 3. The most requested endpoint



No hints. No time limit.

This is what real ops looks like.

```



Patterned directly after real SRE and backend engineering interview tasks.



\---



\## Pricing



| Tier | Price | Access |

|---|---|---|

| \*\*Free\*\* | $0 | Bootcamp + Operation Antariksha (Ch 1-3) + Daily Challenges |

| \*\*Plus\*\* | $12/month | All campaigns + all standalone tracks + Weekly Quests + Interview Mode |

| \*\*Pro\*\* | $24/month | Plus + Team Battles + Private Campaigns + Certificate of Completion |

| \*\*Enterprise\*\* | Custom | Custom campaigns, private leaderboards, SLA, SSO |



\*\*Student discount:\*\* 50% off with `.edu` or `.ac.in` email.



\*\*Team pricing:\*\* Bulk discounts for colleges and companies.



\---



\## Accessibility



LinuxQuest is designed to be accessible:



| Feature | Implementation |

|---|---|

| \*\*Screen reader support\*\* | ARIA labels on all interactive elements, terminal output is accessible |

| \*\*Keyboard navigation\*\* | Full keyboard support (Tab, Enter, Arrow keys) |

| \*\*Color blindness\*\* | ANSI colors are supplemented with symbols (`✅`, `🟡`, `🔴`) |

| \*\*Reduced motion\*\* | Respects OS-level motion preferences |

| \*\*Font sizing\*\* | Terminal font size is adjustable |



The terminal is the primary interface — accessibility here is critical.



\---



\## Tech Architecture



```

linuxquest/

├── frontend/                       # React + Vite

│   ├── components/

│   │   ├── Terminal.jsx                # xterm.js terminal emulator

│   │   ├── StoryPanel.jsx              # Narrative + objectives

│   │   ├── SkillTree.jsx               # Interactive command graph

│   │   ├── HintSystem.jsx              # 3-tier progressive hints

│   │   ├── CommandHandbook.jsx         # Unlockable reference drawer

│   │   ├── Analytics.jsx               # Per-user stats dashboard

│   │   ├── Leaderboard.jsx             # Weekly quest rankings

│   │   └── ProfileCard.jsx             # Shareable Elo + badge card

│   ├── data/

│   │   └── campaigns/

│   │       └── operation-antariksha.js

│   └── hooks/

│       └── useTerminal.js

│

├── backend/                        # Node.js + Express

│   ├── routes/

│   │   ├── execute.js                  # WebSocket → Docker exec

│   │   ├── validate.js                 # Server-side objective checking

│   │   ├── progress.js                 # Save/load user state

│   │   ├── elo.js                      # Rating calculation

│   │   └── weekly.js                   # Weekly quest management

│   ├── sandbox/

│   │   ├── SandboxManager.js           # Abstraction over Docker / Firecracker

│   │   ├── DockerDriver.js             # Current implementation

│   │   ├── FirecrackerDriver.js        # Future: sub-125ms cold starts

│   │   ├── ContainerPool.js            # Pre-warmed idle containers

│   │   └── CommandFilter.js            # Block outbound network in sandbox

│   └── db/

│       └── schema.sql

│

├── sandbox-image/                  # Alpine Linux Docker image

│   ├── Dockerfile

│   ├── setup.sh

│   └── campaigns/

│       └── operation-antariksha/

│           ├── chapter-01/             # Story files per chapter

│           ├── chapter-02/

│           └── ...

│

└── docker-compose.yml

```



\---



\## Core System Design



\### Sandboxed Execution



Every user gets their own ephemeral container. Commands run for real.



```

User types: grep -r "shiva" /var/log

&#x20;       ↓

Frontend → WebSocket → Backend

&#x20;       ↓

SandboxManager.exec(userId, command)

&#x20;       ↓

docker exec <container\_id> grep -r "shiva" /var/log

&#x20;       ↓

Output streams back via WebSocket → xterm.js

```



\*\*Why real containers over simulation:\*\* real output = real learning. Error messages, exit codes, timing, and edge-case behavior all match what players will encounter on a real server.



\*\*Container constraints:\*\*



| Limit | Value |

|---|---|

| CPU | 0.5 core |

| RAM | 128 MB |

| Disk | 512 MB |

| Network | Blocked (Ch9/10: simulated internal net only) |

| Idle timeout | 30 min → auto-destroy |



\*\*Future:\*\* `DockerDriver` and `FirecrackerDriver` both implement `SandboxManager` interface from day one. Firecracker migration = swap one driver, zero architecture changes. Firecracker gives sub-125ms cold starts and stronger kernel-level isolation — critical at 1000+ concurrent users.



\---



\### Objective Validation System



Validators check \*system state\*, not keystrokes. The story advances automatically when state matches. No submit button.



```javascript

{

&#x20; id: "ch2\_obj1",

&#x20; description: "Find all CRITICAL lines in /var/log/isro/anomaly.txt",

&#x20; difficulty: { standard: "grep 'CRITICAL' ...", hard: "grep -E with sort + uniq -c" },

&#x20; hints: \[

&#x20;   "This task involves searching for patterns inside a file",

&#x20;   "The grep command searches text — try: man grep",

&#x20;   "Try: grep 'CRITICAL' /var/log/isro/anomaly.txt"

&#x20; ],

&#x20; validator: {

&#x20;   type: "output\_match",

&#x20;   check: async (container) => {

&#x20;     const r = await exec(container, "grep -c 'CRITICAL' /var/log/isro/anomaly.txt");

&#x20;     return r.stdout.trim() === "14";

&#x20;   }

&#x20; }

}

```



\*\*Validator types:\*\*



| Type | Checks |

|---|---|

| `output\_match` | Command output matches expected |

| `file\_exists` | File created at correct path |

| `file\_content` | File contains specific text |

| `permission\_check` | File has correct chmod bits |

| `process\_dead` | Process no longer running |

| `cron\_absent` | Malicious cron entry removed |

| `env\_var\_set` | Environment variable set correctly |

| `script\_runs` | Script executes without error, produces correct output |

| `network\_blocked` | Domain or IP is unreachable from container |

| `port\_closed` | Previously open port is now closed |

| `package\_installed` | Package present and correct version |

| `service\_state` | systemd unit is in expected state |



All validators run server-side. Players never see the check logic.



\---



\### Elo Rating Engine



```javascript

// Standard Elo with challenge-rating adjustment

function updateElo(playerRating, challengeRating, result, hintsUsed, timeBonus) {

&#x20; const K = 32;

&#x20; const expected = 1 / (1 + 10 \*\* ((challengeRating - playerRating) / 400));

&#x20; const score = result === 'win'

&#x20;   ? Math.max(0.5, 1 - (hintsUsed \* 0.1)) + timeBonus

&#x20;   : 0;

&#x20; return Math.round(playerRating + K \* (score - expected));

}

```



Rating is only affected by Hard Mode completions and Weekly Quest performance. Story mode XP is separate.



\---



\### Story Engine



```javascript

const storyState = {

&#x20; campaign: "operation-antariksha",

&#x20; chapter: 3,

&#x20; objective: 2,

&#x20; hintsUsed: 1,

&#x20; commandHistory: \[],       // full history → used for end-game report + analytics

&#x20; startedAt: Date.now(),

}



// State machine transitions

OBJECTIVE\_COMPLETE  →  show\_story\_beat  →  unlock\_next\_objective

CHAPTER\_COMPLETE    →  chapter\_cinematic  →  unlock\_next\_chapter

ALL\_COMPLETE        →  credits (command history as scrolling code)

```



Story beats: 2–4 lines, typewriter effect. Terminal stays active the whole time.



\---



\### Frontend Layout



```

┌──────────────────────────────────────────────────────────────┐

│  LinuxQuest · Operation Antariksha     \[Ch 3/11] \[📖] \[👤]  │

├─────────────────────────┬────────────────────────────────────┤

│  STORY PANEL            │  TERMINAL                          │

│                         │                                    │

│  Chapter 3: The Hunt    │  arjun@mumbai-node:\~$              │

│                         │  > ps aux                          │

│  Dr. Mehra's voice      │  USER  PID   %CPU  COMMAND         │

│  crackles over comm:    │  root  1     0.0   /sbin/init      │

│  "Arjun, it's on        │  root  847   99.9  shiva\_daemon    │

│   port 4444..."         │  ...                               │

│                         │  arjun@mumbai-node:\~$ \_            │

│  ─────────────────      │                                    │

│  🎯 Kill the rogue      │                                    │

│     process             │                                    │

│                         │                                    │

│  ✅ Find the process    │                                    │

│  ✅ Identify its port   │                                    │

│  ⬜ Kill it             │                                    │

│  ⬜ Verify it's gone    │                                    │

│                         │                                    │

│  \[? Hint]  1/3 used     │                                    │

└─────────────────────────┴────────────────────────────────────┘

```



`\[📖]` → Command Handbook drawer. `\[👤]` → Profile + Analytics.



\---



\### Tech Stack



| Concern | Choice | Reason |

|---|---|---|

| Terminal emulator | `xterm.js` | Industry standard, ANSI, resize, copy-paste |

| Backend | Node.js + Express | `dockerode` SDK, fast WebSocket handling |

| Sandbox (now) | Docker + Alpine Linux | Real Linux, lightweight, chapter snapshots |

| Sandbox (future) | Firecracker MicroVMs | Sub-125ms cold starts, kernel isolation, 1000+ concurrent |

| WebSocket | `socket.io` | Bidirectional streaming for terminal I/O |

| Auth | JWT + Google OAuth | No password management overhead |

| DB | PostgreSQL | Users, progress, command history, Elo ratings |

| Frontend | React + Vite | Fast DX, component model fits campaign engine |

| Deployment | Hetzner VPS + Docker Compose | Cost-effective, full Docker socket control |



\---



\## Build Phases



\### Phase 1 — Playable Demo (4–6 weeks)

\- Chapters 1–3 of Operation Antariksha

\- Docker sandbox + WebSocket terminal

\- 4 validator types: `output\_match`, `file\_exists`, `process\_dead`, `permission\_check`

\- The Bootcamp

\- Session-based (no auth)

\- Command Handbook drawer

\- Deployable at a real URL



\### Phase 2 — Core Platform (4 weeks)

\- Chapters 4–8

\- All 12 validator types

\- 3-tier hint system with XP cost

\- Google auth + progress persistence

\- Daily Challenge system

\- Basic analytics dashboard



\### Phase 3 — Competition Layer (3 weeks)

\- Chapters 9–11 + networking chapter

\- Weekly Quest engine + leaderboard

\- Linux Elo rating

\- Shareable profile cards

\- Real World Mode challenges

\- Shell Scripting Track (6 chapters)



\### Phase 4 — Full Platform

\- All standalone tracks (Git, Docker, Kubernetes, Vim Basics + Mastery, Storage, Packages)

\- Skill Tree visualization (Command Graph)

\- Interview Mode

\- Team Battles

\- Campaign editor + community marketplace

\- Firecracker MicroVM migration

\- Seasonal events

\- Command Line Hero feature



\---



\## Key Risks \& Mitigations



| Risk | Mitigation |

|---|---|

| Docker container escape | `--security-opt=no-new-privileges`, drop all capabilities, read-only root FS + tmpfs overlay |

| Container sprawl | 1 container per user hard limit, 30 min idle → auto-destroy |

| `rm -rf /` inside container | Container is ephemeral — let it happen. Auto-recreate for next objective. |

| Slow container startup | Pre-warmed pool; assign on session start. Firecracker solves this at scale. |

| Cheating (inspecting validator logic) | All validators server-side only, never exposed to client |

| Elo manipulation (alt accounts) | Rate-limit account creation + email verification + IP fingerprinting |

| Community campaign abuse | Review queue for published campaigns; sandboxed validator execution |

| `SandboxManager` abstraction leaking | Docker and Firecracker drivers share interface from day 1; tested independently |

| Vim panic mode overwhelm | Split into Vim Basics (required) and Vim Mastery (optional) |



\---



\## Screenshots



\*(Mockups — to be replaced with real UI)\*



```

\[ Chapter Map ]        \[ Story + Terminal ]    \[ Skill Tree / Command Graph ]



&#x20; ✅ Ch1  100%           ┌──────┬──────────┐     Navigation

&#x20; ✅ Ch2   88%           │Story │ Terminal │     ├── ls   ✅

&#x20; 🔓 Ch3   --            │      │          │     ├── pwd  ✅

&#x20; 🔒 Ch4                 └──────┴──────────┘     ├── find 🟡

&#x20; 🔒 Ch5                                         └── tree 🔴



\[ Weekly Quest ]       \[ Elo Profile ]         \[ Interview Mode ]



&#x20; #17 "Lost Satellite"   @divyansh\_v15           Production 502s.

&#x20; ────────────────        Rating: 1,742           Diagnose it.

&#x20; #1 divyansh  14m22s     Rank: Top 5%            No hints.

&#x20; #2 rahul     16m11s     Streak: 14d             No story.

&#x20; #3 aman      18m53s     Badges: 👑🌐⚡           Real clock.

```



\---



\## Name



\- \*\*LinuxQuest\*\* — the platform

\- \*\*Operation Antariksha\*\* — the flagship campaign

\- \*\*chmod 777: A Love Story\*\* — still the best name, still not using it

```

