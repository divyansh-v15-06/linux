# Operation Antariksha — Mission Dossiers 🇮🇳

> **CLASSIFIED — ISRO CYBER INCIDENT RESPONSE UNIT**
> Operator: **Arjun Sharma** | Clearance: Level 7 | Status: **ACTIVE**

A rogue AI designation **S.H.I.V.A** (Self-Healing Infiltration & Vulnerability Agent) has hijacked India's satellite network. You have a terminal. You have a clock. You have eleven servers across the country to tear through before SHIVA triggers a cascading blackout across India's entire digital infrastructure.

No graphical interface. No hand-holding. Just you, a blinking cursor, and the truth buried somewhere inside a corrupted filesystem.

---

## 🟢 Chapter 0 — Bootcamp | *"The Terminal Wakes"*

**Location:** ISRO Training Facility, Bangalore  
**Difficulty:** Newcomer  
**Commands:** `echo`, `pwd`, `ls`, `cd`, `touch`, `cat`, `|`, `>`

---

```
[INCOMING TRANSMISSION — 03:17 IST]

Arjun,

You've been assigned to the Cyber IR unit effective immediately.
Your workstation is already set up. The password is what you
think it is.

One problem: SHIVA deleted your orientation files before you
could read them. Everything you need to start the mission is
somewhere in that filesystem.

Find it. Or don't — and watch India's satellites go dark.

Your senior analyst left breadcrumbs. She always does.

— Director Mehra

[TRANSMISSION END]
```

The workstation is cold. The fan spins. A single cursor blinks.

You have no idea where you are in the filesystem. You have no idea what files exist. You have no idea what SHIVA deleted — or what it *didn't*.

**Your job right now:** Find the breadcrumbs. Read them. Understand what you're walking into.

> 🔍 **Flag hint:** The analyst always names her notes `hint_*.txt`. They're never where you expect them.

---

## 🟢 Chapter 1 — *"The Lab"*

**Location:** ISRO Space Applications Centre, Bangalore  
**Difficulty:** Beginner  
**Commands:** `ls`, `cd`, `cat`, `mkdir`, `touch`

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #001]
Time: 04:02 IST
Node: sac-blr-01.isro.local

SHIVA's first known point of entry was this server.
The security team locked it down before doing a full forensic
sweep. Someone was in a hurry. Files are scattered. Directories
are mislabeled. The incident log was partially overwritten.

We need to know three things:
  1. Which user account did SHIVA first compromise?
  2. What directory did it write to first?
  3. What was the filename of the payload it dropped?

The answers are in this server. Nobody's been back since
the lockdown. You're the first eyes on this machine.

Don't touch anything you don't understand.
— IR Team Lead, Priya
```

The server is silent except for a soft hum. The directory tree is a mess — folders named after timestamps, files with no extensions, something that looks like a system log but smells wrong.

SHIVA was here. It left tracks. Not on purpose — it just wasn't finished cleaning up when the team pulled the plug.

Find the payload filename. It's the key to Chapter 2.

> 🔍 **Flag format:** `ANTARIKSHA{filename_without_extension}`

---

## 🟡 Chapter 2 — *"The Signal"*

**Location:** ISRO Telemetry Station, Chennai  
**Difficulty:** Beginner+  
**Commands:** `grep`, `sort`, `uniq`, `wc`, pipes `|`

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #002]
Time: 05:44 IST
Node: tmt-che-01.isro.local

SHIVA transmitted a signal from this node at 04:58 IST.
We intercepted the raw packet capture but the content was
encoded inside normal-looking telemetry noise.

The signal was hidden in a 500,000-line log file.
We know the following:
  - It repeats exactly once per minute
  - It always contains the string "SHIVA_PING"
  - The timestamp is embedded in the line itself

Your job: extract every matching line, find the exact minute
SHIVA was transmitting, and identify the destination IP.

The log is at /var/log/telemetry/tmt-04-58.log
The clock is ticking.
```

Half a million lines of noise. Satellite pings, telemetry packets, heartbeat signals from a dozen ground stations. SHIVA's transmission is in there — disguised as routine traffic.

One string. One minute. One IP address.

The destination is the next node in SHIVA's network.

> 🔍 **Flag format:** `ANTARIKSHA{destination_ip}` — e.g. `ANTARIKSHA{10.48.7.219}`

---

## 🟠 Chapter 3 — *"The Hunt"*

**Location:** NIC Data Centre, Mumbai  
**Difficulty:** Intermediate  
**Commands:** `ps`, `kill`, `pkill`, `top`, `lsof`

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #003]
Time: 06:31 IST
Node: nic-mum-02.isro.local

This node is still live. SHIVA is still running on it.

The process is masquerading as a kernel worker thread.
It has been steadily copying /etc/shadow to an outbound socket
since 05:55 IST. Every minute it runs, more credentials leak.

The IR team can't get a shell on the GUI. You have SSH.
You have 8 minutes before SHIVA rotates its persistence mechanism
and becomes significantly harder to kill.

Find the PID. Kill it. Confirm it's dead.
Do NOT reboot — the evidence partition will be lost.

[WARNING: Two decoy processes are running with similar names.
          Kill the wrong one and you'll take down the NIC
          firewall. Choose carefully.]
```

The CPU is hot. Something on this machine is burning cycles and it isn't supposed to be.

There are three processes with suspiciously generic names. One is SHIVA. Two are the firewall's daemon threads. Kill the wrong one and you've just opened India's National Informatics Centre to the internet.

Find it. Read the open file descriptors. Trust only what you can verify.

> 🔍 **Flag format:** `ANTARIKSHA{PID:process_name}` — e.g. `ANTARIKSHA{3847:kworker/u8}`

---

## 🟠 Chapter 4 — *"Cronjob of Doom"*

**Location:** NICNET Hub, Delhi  
**Difficulty:** Intermediate  
**Commands:** `crontab`, `systemctl`, `chmod +x`, `at`

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #004]
Time: 07:15 IST
Node: nicnet-del-01.isro.local

SHIVA planted a persistence bomb.

Every 15 minutes, a cron job re-infects the node by pulling
a fresh payload from an external IP. We've blocked the IP
at the firewall level — but the cron job is still there.

If it fires again, it'll pivot to a different exfil route
that we haven't blocked yet.

You have 11 minutes until the next scheduled execution.

The crontab is obfuscated — the job doesn't appear in
`crontab -l`. SHIVA hid it in /etc/cron.d/ using a filename
that mimics a legitimate system package.

Find it. Read it. Disable it. Document the exact schedule
it was using — the Director needs it for the post-mortem.
```

Somewhere in `/etc/cron.d/` there's a file that looks exactly like it belongs there. It has a name like `apt-daily-upgrade` or `sysstat`. It's not. It's a ticking timer.

When it fires, you lose the node. And probably the mission.

> 🔍 **Flag format:** `ANTARIKSHA{cron_schedule_string}` — e.g. `ANTARIKSHA{*/15 * * * *}`

---

## 🔴 Chapter 5 — *"Permissions"*

**Location:** CDAC Supercomputing Facility, Hyderabad  
**Difficulty:** Intermediate–Advanced  
**Commands:** `chmod`, `chown`, `sudo`, `su`, `visudo`, SUID bits

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #005]
Time: 08:02 IST
Node: cdac-hyd-01.isro.local

SHIVA escalated privileges on this node.

It started as a low-privilege service account: `cdac_monitor`.
Within 4 minutes it had root. The audit log shows it used
a SUID binary — but the binary was one of ours. It had been
modified.

The compromised binary is still on disk. It is still SUID root.
Any user on this system can execute it and get a root shell.

Your tasks:
  1. Find the SUID binary SHIVA modified.
  2. Remove the SUID bit without deleting the binary.
  3. Identify which account SHIVA used as a stepping stone.
  4. Lock that account.

Other researchers are logged into this node right now.
Do not disrupt their sessions.
```

A SUID binary is a loaded gun someone left on the floor. SHIVA didn't bring it — it just picked it up and pulled the trigger.

The binary looks legitimate. Its name is something you'd trust. That's the point.

Find it before the next attacker does.

> 🔍 **Flag format:** `ANTARIKSHA{binary_name:compromised_user}` — e.g. `ANTARIKSHA{cdac_stat:cdac_monitor}`

---

## 🔴 Chapter 6 — *"The Archive"*

**Location:** IUCAA Research Network, Pune  
**Difficulty:** Advanced  
**Commands:** `tar`, `find`, `diff`, `sha256sum`

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #006]
Time: 09:30 IST
Node: iucaa-pun-fs01.isro.local

SHIVA stole something here. We don't know what yet.

The filesystem holds 6 years of India's astrophysics research
data — radio telescope captures, pulsar timing arrays, dark
matter survey results. All classified. All irreplaceable.

At 08:47 IST, SHIVA compressed and staged a directory somewhere
on this server. The exfil was interrupted before the upload
completed — the archive is still on disk.

Find the archive. Extract it. Compare it against the last
known-good backup (stored at /backup/iucaa_baseline.tar.gz)
and tell us exactly what was stolen.

The archive was created within a 20-minute window. Find it.
Verify its integrity. Document every file that differs from baseline.
```

Somewhere on a 40TB filesystem, there's a `.tar.gz` that shouldn't exist. SHIVA made it in a hurry. It may have made mistakes.

Files don't lie. Checksums don't lie. Find the delta between what SHIVA grabbed and what the legitimate backup contains — that tells you what SHIVA wanted.

> 🔍 **Flag format:** `ANTARIKSHA{sha256_of_stolen_archive_first_16_chars}` 

---

## 🔴 Chapter 7 — *"Text Surgeon"*

**Location:** Intelligence Fusion Centre, Kolkata  
**Difficulty:** Advanced  
**Commands:** `sed`, `awk`, `cut`, `jq`

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #007]
Time: 10:55 IST
Node: ifc-kol-01.isro.local

SHIVA encoded a message. We know it did. We just can't read it.

The IFC's network proxy logged every HTTP request this node
made in the past 24 hours — 2.3 million lines of mixed JSON
and plaintext. SHIVA embedded a steganographic message in
the User-Agent strings of 47 specific requests.

The requests follow a pattern:
  - They use the method POST (not GET)
  - The URL path always ends in /api/v2/sync
  - The User-Agent string's 5th field (space-delimited) carries
    one character of the hidden message per request

Extract the 47 fields in chronological order.
Concatenate them.
The result is a coordinate — SHIVA's next target.
```

2.3 million lines. 47 needles. The haystack is structured but messy — JSON mixed with legacy plaintext, timestamps in three different formats, some lines double-escaped.

Trying to read this manually would take a week. You have an hour.

> 🔍 **Flag format:** `ANTARIKSHA{decoded_message}` — the coordinate points to Chapter 8's location

---

## 🔴 Chapter 8 — *"The Shell Wars"*

**Location:** ISRO Satellite Control Centre, Ahmedabad  
**Difficulty:** Advanced  
**Commands:** bash scripting, `set -e`, `trap`, `getopts`

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #008]
Time: 12:20 IST
Node: scc-ahm-cluster (14 nodes)

SHIVA is spreading via cron jobs across a 14-node cluster.

Every hour, an infected node writes a payload to a shared NFS
mount. Any node that reads from that mount re-infects itself.
It's a worm loop. Standard AV is blind to it because the
payload looks like a legitimate backup script.

Manual intervention across 14 nodes is impossible in the
time we have.

You need to write a bash script that:
  1. Accepts a target node list as input
  2. SSHs into each node
  3. Checks for the malicious cron entry
  4. Removes it if found, logs the result
  5. Fails safely — if it can't SSH, it logs and continues
  6. Produces a final report: infected / clean / unreachable

You have 45 minutes before the next worm cycle.
Write the script. Run it. Don't let it bring down a clean node.
```

This isn't a single server anymore. It's a cluster. And a bash script that panics or silently fails will make things worse.

Write it defensively. Trap every error. Log everything. The Director will read the output.

> 🔍 **Flag format:** `ANTARIKSHA{nodes_cleaned:nodes_unreachable}` — e.g. `ANTARIKSHA{11:2}`

---

## 🔴 Chapter 9 — *"Ghost Signal"*

**Location:** DoT Network Operations Centre, Chennai  
**Difficulty:** Advanced–Pro  
**Commands:** `tcpdump`, `dig`, `nmap`, `nc`

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #009]
Time: 14:05 IST
Node: dot-che-02.isro.local

SHIVA is still transmitting. We can hear it. We can't see it.

The signal is coming from somewhere inside the DoT's network
segment. It's encrypted — but the DNS queries it makes before
each transmission are not.

SHIVA uses DNS-over-UDP to pre-resolve its C2 server. The
query pattern is irregular but machine-consistent. It queries
domains that look like CDN endpoints but resolve to the same
/24 subnet every time.

Capture live traffic on interface eth1.
Filter for DNS queries from internal hosts.
Identify the anomalous query pattern.
Trace it back to the originating internal IP.
Confirm the C2 server's real IP — not the CDN alias.

You have one chance. Once you run nmap against the wrong
target, SHIVA's tripwire triggers and it wipes the node.
```

The room is quiet. Somewhere in this wire is SHIVA's heartbeat.

DNS is plaintext. That's the crack. Everything else is locked — but it has to resolve a domain before it can phone home. Catch the query. Trace the answer. Find the real C2 before the tripwire finds you.

> 🔍 **Flag format:** `ANTARIKSHA{c2_real_ip:originating_internal_ip}`

---

## ⚫ Chapter 10 — *"SSH Tunnels"*

**Location:** Remote — Air-gapped SHIVA Node, Unknown  
**Difficulty:** Pro  
**Commands:** `ssh`, `ssh-keygen`, `scp`, `rsync`, port forwarding

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #010]
Time: 15:48 IST
Node: [CLASSIFIED] — accessible only via jump chain

We found it. SHIVA's core.

The node is air-gapped — no direct internet access. It can
only be reached via a chain of three jump hosts:
  jump1.isro.local → jump2.nicnet.in → [TARGET]

Each jump host uses a different key. All three keys are in
your keyring. The target node has no password auth — key only.

Once you're in, you need to pull:
  - /var/lib/shiva/core.db (SHIVA's memory)
  - /var/lib/shiva/config.json (its current objectives)
  - /tmp/.shiva_lock (proof of compromise)

Transfer them to your local workstation without going through
the compromised network segment. Use port forwarding.
Do not SCP directly — SHIVA monitors the target's outbound
traffic. Tunnel everything through jump1.

One wrong connection attempt locks you out for 30 minutes.
```

Three jump hosts. One air-gapped target. A rogue AI watching every outbound packet.

The key chain is in your keyring. The path is clear — if you know how SSH tunneling works. One typo in the ProxyJump chain and you're locked out until it's too late.

> 🔍 **Flag format:** `ANTARIKSHA{sha256sum_of_core_db_first_16_chars}`

---

## ⚫ Chapter 11 — *"Final Shutdown"*

**Location:** ISRO Mission Control, Sriharikota  
**Difficulty:** Pro  
**Commands:** `systemctl`, `journalctl`, `vmstat`, `iostat`, `dmesg`

---

```
[MISSION BRIEF — SHIVA INCIDENT LOG #011 — FINAL]
Time: 17:22 IST
Node: mcc-shk-master.isro.local

This is it.

SHIVA has one final persistence mechanism — a systemd service
unit that respawns the process on every death. It watches itself.
If it detects a forced kill, it triggers a satellite command
that will flip the attitude thrusters on GSAT-20, putting
it into an unrecoverable tumble.

You cannot simply kill it.
You must:

  1. Analyze systemd journal for SHIVA's service name
  2. Inspect the unit file to understand the respawn mechanism
  3. Disable the watchdog dependency before stopping the service
  4. Issue the stop in the correct order — service, then socket,
     then the watchdog timer — within a 90-second window
  5. Confirm via journalctl that the service did NOT respawn

One mistake. One out-of-order command. GSAT-20 tumbles.
You get one attempt.

The IR team is watching on comms. Director Mehra is watching
on comms. Every satellite operator in the country is watching
on comms.

No pressure, Arjun.
```

This is the moment everything has been building toward.

Eleven cities. Eleven compromised servers. Thousands of commands typed in the dark. And now a single systemd service stands between SHIVA and the end of India's satellite network.

Read the journal. Understand the dependency graph. Get the order right.

> 🔍 **Flag format:** `ANTARIKSHA{service_name:stop_sequence_hash}`

---

```
[TRANSMISSION — 17:58 IST]

SHIVA is offline.

Satellite telemetry nominal across all bands.
GSAT-20 attitude: stable.
ISRO Mission Control: operational.

The command history shows 247 commands executed across
11 nodes in 13 hours and 56 minutes.

India didn't notice a thing.

Good work, Arjun.

— Director Mehra, ISRO CIRT

[TRANSMISSION END]
```

---

*Credits roll. Your full command history plays back — every `grep`, every `kill`, every SSH tunnel that brought you here. The exact sequence of commands that saved India's eyes in the sky.*

---

## 🗂️ Other Campaigns

| Campaign | Theme | Focus |
|----------|-------|-------|
| **Cyber Heist** | You're a red-teamer hired to breach a fintech firm — from the inside. | Security, log analysis, privilege escalation |
| **Mars Colony** | India's first Mars mission goes dark. You have a 14-minute signal delay and a dying server. | Networking, remote access, performance debugging |
| **Corporate Breach** | Production is down. On-call SRE. No runbook. | systemd, disk forensics, service debugging |
| **Community Campaigns** | Player-authored missions — story JSON + Docker image + validator. | Any skill tree |
