# LinuxQuest — Operation Antariksha 🇮🇳

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Status: Concept](https://img.shields.io/badge/Status-Concept-orange)](#)
[![Campaign: Operation Antariksha](https://img.shields.io/badge/Campaign-Operation%20Antariksha-blue)](#operation-antariksha)

> *A rogue AI has hijacked India's satellite network. You have a terminal. You have a clock. Eleven servers across the country are waiting.*

---

## What is LinuxQuest?

**LinuxQuest** is a CTF-style, story-driven Linux learning platform. Players learn real Linux administration — not through tutorials, but through mystery missions set across India, where every command is earned by needing it.

You don't study `kill`. You use it to stop a rogue process before it wipes a server.

**Learn by doing. Under pressure. In the dark.**

---

## The Flagship Campaign — Operation Antariksha

You are **Arjun Sharma**, a junior ISRO engineer on your first week. A rogue AI designation **S.H.I.V.A** (Self-Healing Infiltration & Vulnerability Agent) has hijacked India's satellite network.

No GUI. No hints. Just a terminal and eleven compromised servers across the country.

> Full mission dossiers with classified briefings, flag formats, and story → [`story.md`](./story.md)

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

**Ending:** SHIVA goes offline. The satellite network comes back. Credits roll with your full command history — the exact sequence of commands that saved India.

---

## Other Campaigns

| Campaign | Premise | Focus |
|----------|---------|-------|
| **Cyber Heist** | Red-teamer hired to breach a fintech firm from the inside | Security, log analysis, privilege escalation |
| **Mars Colony** | India's first Mars mission goes dark — 14-min signal delay, dying server | Networking, remote access, performance |
| **Corporate Breach** | Production is down. On-call SRE. No runbook. | systemd, disk forensics, service debugging |
| **Community Campaigns** | Player-authored missions — story JSON + Docker image + validator | Any skill tree |

---

## Platform Features

- **Skill Trees** — Track mastery across Core Linux, Networking, Security, DevOps, and Performance. Green = mastered. Red = weak.
- **Daily Challenges** — A single mystery task each day. 5–10 minutes. Streak system like Duolingo.
- **Weekly Quests** — Time-limited, competitive missions. Same clean container for everyone. Global leaderboard.
- **Linux Elo Rating** — Global skill rating (Newcomer → Wizard). Shareable. Verifiable via API.
- **Interview Mode** — Real server tasks. No story. No hints. Timed.
- **Team Battles** — College vs college. Club competitions. Hackathons.

---

## Difficulty Levels

Every objective ships in two modes. Players choose before starting.

| Mode | Example |
|------|---------|
| **Standard** | `grep "ERROR" anomaly.txt` |
| **Hard** | `grep -E "ERROR\|CRITICAL" /var/log/*.log \| sort \| uniq -c \| sort -rn \| head -20` |

Hard mode: 2× XP, exclusive badges, counts toward Elo rating.

---

## Tracks (Standalone)

| Track | Story Hook | Commands |
|-------|-----------|----------|
| Shell Scripting | SHIVA is spreading via cron jobs across 100+ nodes | `bash`, `trap`, `set -e`, `getopts` |
| Storage & Filesystems | Corrupted satellite backup disk — recover the telemetry | `df`, `mount`, `fsck`, `fdisk` |
| Package Management | Forensic toolkit missing. 8 minutes before SHIVA rotates keys | `apt`, `dnf`, `rpm`, `dpkg` |
| Vim Basics | Remote server. No nano. 15-minute window to edit the launch config | `vim` — normal/insert/save |
| Vim Mastery | SSH drops every 30 seconds. Fix the config before the satellite launch window closes | `vim` — macros, registers, visual mode |
| Git | Developer force-pushed to main. 3 weeks of satellite firmware lost. Recover it | `git reflog`, `cherry-pick`, `bisect` |
| Docker | SHIVA escaped into Docker containers. Find it. Kill it. Prevent it restarting | `docker ps`, `exec`, `logs`, `network` |
| Kubernetes | SHIVA replicated across a live cluster. Isolate every instance | `kubectl get/logs/exec/taint` |
| Environment Variables | Launch codes stored as env vars on a locked system. Recover without rebooting | `export`, `env`, `printenv`, `.env` |

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

Every player has a global Linux rating — like Chess Elo.

```
< 800     Newcomer
800–1199  User
1200–1599 Power User
1600–1999 Sysadmin
2000–2399 Engineer
2400+     Wizard
```

Recruiter-shareable. Embeddable on GitHub README or LinkedIn.

---

## Roadmap

| Phase | Goal |
|-------|------|
| **Phase 1 — MVP** | Bootcamp + Operation Antariksha (Ch 1–3) + Validator Engine |
| **Phase 2 — Competition Layer** | Weekly Quests + Elo Rating + Leaderboard |
| **Phase 3 — Community Layer** | Community Campaigns + Tracks + Profile Cards |
| **Phase 4 — Enterprise Layer** | Interview Mode + Team Battles + API |

---

## Contributing

See [`CONTRIBUTING.md`](./CONTRIBUTING.md) for campaign format, validator spec, coding standards, and how to submit new tracks.

---

## License

MIT © 2026–2027 Divyansh.

---

*Full mission dossiers, classified briefings, and flag formats → [`story.md`](./story.md)*
