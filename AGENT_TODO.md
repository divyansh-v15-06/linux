# LinuxQuest — Agent TODO

> **For AI Agents:** This file is your source of truth. Read it fully before taking any action.
> Mark tasks `[/]` when starting, `[x]` when done. Never skip acceptance criteria.
> Always read `## Context` and `## Constraints` before writing code.

---

## Context

**Project:** LinuxQuest — a CTF-style, story-driven Linux learning platform.
**Flagship campaign:** Operation Antariksha (11 chapters, story in `story.md`).
**Stack:**
- Frontend: React + Vite (TypeScript) — single-page, shell-only UI
- Backend: Go + Gin (Fly.io free tier)
- Database: PostgreSQL (Supabase free tier)
- Auth: Google OAuth 2.0 + JWT
- Email: Gmail SMTP with App Password
- **Sandbox: CheerpX (WebAssembly Linux — runs in the browser, zero server cost)**
- Disk images: Cloudflare R2 free tier (10GB)
- Realtime: WebSockets (leaderboard updates only — no terminal I/O to server)

**Key files to read before coding:**
- `README.md` — full platform design (merged, authoritative)
- `story.md` — all 11 chapter dossiers, flag formats, narrative
- `FRONTEND.md` — shell-only UI spec, command map, virtual filesystem

---

## Constraints

- Do NOT skip phases — each phase depends on the previous.
- Do NOT build a traditional UI — the frontend is a shell. No nav bars, no buttons, no modals. Read `FRONTEND.md`.
- Do NOT expose flag values or validator logic to the client — server-side only.
- Do NOT use placeholder data in UI — pull content from `story.md`.
- Do NOT use Docker or any server-side sandbox — Linux runs in the browser via CheerpX.
- Do NOT expose flag values or validator logic to the client — flag hashes compared server-side only.
- Every new feature needs a working test before marking `[x]`.
- Disk images are pre-built and hosted on Cloudflare R2 — never generated at runtime.

---

## Phase 1 — Foundation & Repo Setup

### 1.1 Project Scaffold
- [ ] Initialize monorepo structure:
  ```
  linuxquest/
  ├── app/          ← React + Vite frontend (shell UI)
  ├── server/       ← Go + Gin backend (Fly.io)
  ├── images/       ← Linux disk images per chapter (uploaded to R2)
  ├── campaigns/    ← JSON campaign + chapter definitions
  ├── docs/         ← Architecture, API docs
  └── scripts/      ← Dev utilities, seed scripts, image builder
  ```
- [ ] Add root `Makefile` with targets: `dev`, `build`, `test`, `lint`, `sandbox`
- [ ] Add `.env.example` with all required env vars documented:
  - `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `GOOGLE_REDIRECT_URI`
  - `GMAIL_USER`, `GMAIL_APP_PASSWORD`
  - `DATABASE_URL`
  - `JWT_SECRET`
- [ ] Add `docker-compose.yml` for local dev (postgres, backend, frontend, sandbox)
- [ ] Add GitHub Actions CI: lint + test on every PR

**Acceptance criteria:** `make dev` starts the full stack locally. `make test` passes with 0 failures.

---

### 1.2 Database Schema
- [ ] Design and migrate initial schema:
  ```sql
  users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    google_id TEXT UNIQUE,
    elo INT DEFAULT 800,
    xp INT DEFAULT 0,
    level INT DEFAULT 1,
    streak INT DEFAULT 0,
    last_active TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
  )

  chapters (
    id UUID PRIMARY KEY,
    campaign_slug TEXT,
    number INT,
    title TEXT,
    city TEXT,
    difficulty TEXT,
    commands TEXT[],
    flag_hash TEXT,    ← SHA-256 of the flag, never plaintext
    story_text TEXT    ← from story.md
  )

  user_progress (
    user_id UUID REFERENCES users,
    chapter_id UUID REFERENCES chapters,
    status TEXT CHECK (status IN ('locked','active','complete')),
    attempts INT DEFAULT 0,
    hints_used INT DEFAULT 0,
    completed_at TIMESTAMP,
    PRIMARY KEY (user_id, chapter_id)
  )

  submissions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users,
    chapter_id UUID REFERENCES chapters,
    flag_submitted TEXT,
    is_correct BOOLEAN,
    submitted_at TIMESTAMP DEFAULT NOW()
  )

  sandbox_sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users,
    chapter_id UUID REFERENCES chapters,
    container_id TEXT,
    started_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
  )

  command_history (
    id UUID PRIMARY KEY,
    user_id UUID,
    chapter_id UUID,
    command TEXT,
    ran_at TIMESTAMP DEFAULT NOW()
  )
  ```
- [ ] Write seed data for Operation Antariksha (Ch 0–11) from `story.md` — including flag hashes
- [ ] Write migration scripts (up + down)

**Acceptance criteria:** `make db-migrate` runs clean. Seed populates all 12 chapters with hashed flags.

---

## Phase 2 — Backend API

### 2.1 Auth — Google OAuth + JWT
- [ ] `GET /api/auth/google` — redirect to Google OAuth consent screen
- [ ] `GET /api/auth/google/callback` — receive code, exchange for tokens, get user info
- [ ] On first login: create user in DB, send welcome email via Gmail SMTP
- [ ] On any login: issue JWT (15 min) + refresh token (7 days), return to frontend via `postMessage`
- [ ] `POST /api/auth/refresh` — exchange refresh token for new JWT
- [ ] `POST /api/auth/logout` — revoke refresh token
- [ ] JWT middleware — attach user to context on all protected routes
- [ ] Rate limiting on auth endpoints (max 5 req/min per IP)

**Gmail SMTP setup (Go):**
```go
auth := smtp.PlainAuth("", os.Getenv("GMAIL_USER"), os.Getenv("GMAIL_APP_PASSWORD"), "smtp.gmail.com")
smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, body)
```

**Acceptance criteria:** Full OAuth flow works end-to-end. Welcome email sent on first login. Invalid/expired JWTs return 401.

---

### 2.2 Campaign & Chapter API
- [ ] `GET /api/campaigns` — list campaigns with status for authenticated user
- [ ] `GET /api/campaigns/:slug` — campaign detail + chapter list with user progress
- [ ] `GET /api/campaigns/:slug/chapters/:num` — chapter detail (story, commands, difficulty) — **no flag in response**
- [ ] `POST /api/campaigns/:slug/chapters/:num/submit` — validate flag (compare SHA-256 hash), update progress, award XP
- [ ] `GET /api/users/:username/progress` — all chapter statuses for a user
- [ ] `GET /api/users/:username/profile` — full profile (XP, Elo, badges, streak, skill breakdown)

**Acceptance criteria:** Flag submission correctly validates correct flags, rejects wrong ones. Flag value never appears in any API response.

---

### 2.3 Progression Gate
- [ ] Backend enforces strict linear progression:
  ```go
  func canAccessChapter(userID uuid.UUID, chapterNum int) bool {
      if chapterNum == 0 { return true } // bootcamp always open
      prev := db.GetChapterStatus(userID, chapterNum - 1)
      return prev == "complete"
  }
  ```
- [ ] `cd ch<N>` on locked chapter → backend returns 403 with locked message
- [ ] Shell displays: `[ACCESS DENIED] Complete Ch N first. Progress: ...`
- [ ] Exception: `skip bootcamp` → runs 5-command validation → marks Ch 0 complete if passed
- [ ] Completed chapters are always re-enterable (replay mode, no XP re-awarded)

**Acceptance criteria:** User cannot access Ch 5 without completing Ch 4. Replay of completed chapters works without resetting progress.

---

### 2.4 Sandbox — CheerpX Image API

> **No Docker. No VPS for sandbox. Linux runs in the user's browser via CheerpX (WebAssembly).**
> This section is only about serving disk images — not running containers server-side.

- [ ] `GET /api/sandbox/image-url/:chapterNum` — return a signed Cloudflare R2 URL for the chapter disk image
  - Check progression gate before returning URL
  - URL expires in 15 minutes (player must be logged in and chapter unlocked)
- [ ] Cloudflare R2 bucket setup:
  - Bucket: `linuxquest-images` (free tier: 10GB storage, 10M reads/month)
  - One `.img` file per chapter: `ch0.img`, `ch1.img` … `ch11.img`
  - Images are read-only and pre-built via `scripts/build-image.sh`

**Acceptance criteria:** API returns a valid R2 URL only if the chapter is unlocked. URL is time-limited. No server-side process is spawned.

---

### 2.5 Leaderboard & Elo API
- [ ] `GET /api/leaderboard/weekly` — top 50 by weekly XP
- [ ] `GET /api/leaderboard/global` — top 50 by Elo
- [ ] `GET /api/users/:username/elo` — Elo + rank band
- [ ] Elo update on Hard mode chapter completion and Weekly Quest:
  ```go
  func updateElo(playerRating, challengeRating int, result string, hintsUsed int) int {
      K := 32.0
      expected := 1.0 / (1.0 + math.Pow(10, float64(challengeRating-playerRating)/400))
      score := 0.0
      if result == "win" {
          score = math.Max(0.5, 1.0-float64(hintsUsed)*0.1)
      }
      return playerRating + int(math.Round(K*(score-expected)))
  }
  ```
- [ ] Rating bands: Newcomer (<800), User (800–1199), Power User (1200–1599), Sysadmin (1600–1999), Engineer (2000–2399), Wizard (2400+)
- [ ] Story mode XP is separate from Elo — Elo only affected by Hard mode + Weekly Quests

**Acceptance criteria:** Completing a Hard chapter increases Elo more than Standard. Using hints reduces Elo gain.

---

### 2.6 Email Notifications (Gmail SMTP)
- [ ] Welcome email on first login
- [ ] Weekly quest reminder (Friday 6 PM IST)
- [ ] Badge earned notification
- [ ] Streak at risk warning (if no activity by 9 PM, remind at 8 PM)
- [ ] All emails must be plain-text with optional HTML fallback — no heavy templates
- [ ] Rate limit: max 2 emails/day per user

**Acceptance criteria:** Welcome email arrives within 60 seconds of first login. Emails don't go to spam (verify SPF/DKIM for Gmail SMTP).

---

## Phase 3 — Frontend (Shell-Only UI)

> **Read `FRONTEND.md` in full before writing any frontend code.**
> The entire UI is a shell. No pages, no nav bars, no buttons, no modals.
> Everything happens inside a single `xterm.js` terminal instance.

### 3.1 Shell Setup
- [ ] Single-page React + Vite app — one `<App />` component, no React Router
- [ ] Mount `xterm.js` full-screen on load with `xterm-addon-fit`
- [ ] On mount: play boot sequence (typing animation, ISRO-CIRT banner)
- [ ] Blinking block cursor `█`, font `JetBrains Mono 14px`
- [ ] Color palette (CSS variables only):
  - `--bg: #0d1117` · `--text: #c9d1d9` · `--green: #39ff14`
  - `--red: #ff4444` · `--amber: #ffaa00` · `--dim: #444c56`
- [ ] Thin single-line status bar at top: `[ISRO-CIRT] <cwd> | <XP>` — updates on navigation
- [ ] **Copy-paste disabled** — block all paste vectors:
  - `Ctrl+V` / `Ctrl+Shift+V` via `attachCustomKeyEventHandler` → returns `false`
  - Right-click via `contextmenu` → `preventDefault()`
  - Browser paste event → `preventDefault()` + `stopPropagation()`
  - On paste attempt: print `[PASTE DISABLED] Type the command. No shortcuts here.` in amber

**Acceptance criteria:** Site loads as full-screen black terminal. Boot sequence plays. Prompt appears. No other UI visible. Ctrl+V, right-click paste, and middle-click paste all show the disabled message — nothing is pasted.

---

### 3.2 Virtual Filesystem & Command Parser
- [ ] Implement in-memory virtual filesystem (JS object tree) matching `FRONTEND.md` structure
- [ ] Command parser — reads stdin, dispatches to handlers:
  - `login` — opens Google OAuth popup, resumes shell on success
  - `ls` / `ls -la` — list current directory
  - `cd <dir>` / `cd ..` — navigate virtual FS
  - `pwd` — print current path
  - `cat <file>` — render file contents to terminal
  - `clear` — clear terminal buffer
  - `whoami` — username + Elo band
  - `history` — recent command history
  - `help` — full command reference
  - `start` — inside chapter dir: launch sandbox WebSocket session
  - `submit <flag>` — inside chapter dir: validate flag via API
  - `hint` — inside chapter dir: fetch hint (costs Elo)
  - `skip bootcamp` — runs 5-command validation, marks Ch 0 complete if passed
- [ ] Unknown command → `bash: <cmd>: command not found`
- [ ] After 3 consecutive unknown commands → print amber hint line
- [ ] Tab autocomplete for directory names and known commands
- [ ] `login @email` syntax opens Google OAuth directly

**Acceptance criteria:** Player navigates `cd missions → cd antariksha → cd ch1 → cat brief → start` entirely via shell. No mouse required at any step.

---

### 3.3 Google OAuth Shell Flow
- [ ] `login` command triggers `window.open(googleOAuthURL)` — popup
- [ ] Backend returns JWT via `postMessage` to opener window after OAuth completes
- [ ] Shell receives `postMessage`, updates session state, changes prompt to `arjun@linuxquest`
- [ ] `whoami` returns `arjun — Elo: 1,247 — Power User`
- [ ] `logout` clears JWT, resets to `guest@linuxquest`
- [ ] If accessing a protected command while `guest` → `[AUTH REQUIRED] Type  login  to authenticate`

**Acceptance criteria:** Full login flow works without leaving the terminal. Popup opens, user authenticates, prompt updates — no page reload.

---

### 3.4 Chapter Entry Flow
- [ ] `cd ch<N>` checks progression gate via API:
  - Unlocked: transmission animation → brief auto-prints from `story.md`
  - Locked: `[ACCESS DENIED]` message with progress indicator
- [ ] Prompt hostname changes to chapter server (e.g. `arjun@sac-blr-01`)
- [ ] `start` flow:
  1. Call `GET /api/sandbox/image-url/:chapterNum` → get signed R2 URL
  2. Load CheerpX with the disk image URL: `CheerpXEnv.create({ mounts: [{ type: 'disk', url }] })`
  3. Boot Alpine Linux inside the browser tab — no WebSocket to server needed
  4. Connect CheerpX stdin/stdout to xterm.js directly in the browser
- [ ] `exit` / `Ctrl+D` destroys the CheerpX instance, returns to mission shell
- [ ] `submit ANTARIKSHA{...}` calls submit API, prints green (correct) or red (wrong) result
- [ ] Flag validation still happens server-side — compare SHA-256 hash, never expose plaintext

**Acceptance criteria:** Full chapter loop works terminal-only. Locked chapter shows correct message. Flag submission gives immediate visual feedback.

---

### 3.5 Cat Outputs (rendered inside terminal)
- [ ] `cat help` — full command table, column-aligned
- [ ] `cat profile` — XP ASCII bar (`████░░`), Elo, level, badges, streak
- [ ] `cat leaderboard` — top-20 table with rank, username, Elo, weekly XP
- [ ] `cat daily` — today's challenge brief + objectives
- [ ] `cat missions/README` — campaign list with completion status icons
- [ ] `cat missions/antariksha/README` — ASCII India map with chapter city nodes (locked = `?`, active = blinking, complete = `✓`)
- [ ] `man shiva` / `man arjun` — lore entries in `man` page format
- [ ] Easter eggs: `grep flag *`, `ping shiva`, `ssh shiva@satellite`, `chmod 777 shiva`

**Acceptance criteria:** Every `cat` output fits terminal width without horizontal scrolling. All outputs load in <500ms.

---

### 3.6 CheerpX Browser Integration
- [ ] Install CheerpX: `npm install @leaningtech/cheerpx`
- [ ] On `start` command:
  ```js
  const imageUrl = await api.getSandboxImageUrl(chapterNum); // signed R2 URL
  const cx = await CheerpXEnv.create({
    mounts: [{ type: 'disk', url: imageUrl, passphrase: null }]
  });
  const proc = await cx.spawn('/bin/sh', [], {
    stdin: xterm,
    stdout: xterm,
    stderr: xterm
  });
  ```
- [ ] xterm.js connects directly to CheerpX process stdin/stdout — no server in the loop
- [ ] On `exit`: call `cx.terminate()`, return to mission shell
- [ ] Linux boots in browser in <5s (disk image is downloaded from R2 + cached)
- [ ] Browser caches the disk image after first load — subsequent chapter replays are instant
- [ ] `rm -rf /` is harmless — browser memory only, resets on next `start`
- [ ] Network inside CheerpX: disabled by default (Ch 9–10 use simulated file-based network puzzles instead of live tcpdump)

**Acceptance criteria:** Player types real Linux commands in browser, they execute locally via CheerpX WebAssembly, output appears in <100ms. Zero server load from terminal I/O.

---

## Phase 4 — Disk Image Content (Chapter Filesystems)

> **No Dockerfiles. Build Alpine Linux ext4 disk images using `scripts/build-image.sh`.**
> Each image is uploaded to Cloudflare R2. CheerpX mounts it in the browser.

For each chapter, create `images/chapter-<N>/rootfs/` with pre-seeded filesystem content, then build to `ch<N>.img`:

- [ ] **Ch 0:** Basic Alpine Linux + `hint_*.txt` breadcrumbs hidden in subdirectories
- [ ] **Ch 1:** Scattered files, mislabeled directories, partially overwritten incident log with SHIVA payload filename hidden inside
- [ ] **Ch 2:** 500,000-line telemetry log at `/var/log/telemetry/tmt-04-58.log` with `SHIVA_PING` in exactly 60 lines
- [ ] **Ch 3:** Init scripts that launch 3 processes on boot — one high-CPU SHIVA process, two legitimate daemon lookalikes
- [ ] **Ch 4:** Malicious file in `/etc/cron.d/` mimicking a legitimate package name
- [ ] **Ch 5:** Modified SUID binary in `/usr/bin/` that spawns root shell
- [ ] **Ch 6:** Staged `.tar.gz` at unexpected path + `/backup/baseline.tar.gz` for comparison
- [ ] **Ch 7:** 2.3M-line mixed JSON/plaintext log at `/var/log/proxy/access.log` with 47 steganographic POST entries
- [ ] **Ch 8:** Shell worm script + "shared" directory simulating NFS mount infection across fake node dirs
- [ ] **Ch 9:** Pre-recorded DNS query log + packet capture file (`.pcap`) — player uses `tcpdump -r` not live capture
- [ ] **Ch 10:** Multi-hop SSH config pre-loaded with keys in `/home/player/.ssh/` — `ssh-agent` pre-configured
- [ ] **Ch 11:** systemd service unit files pre-installed with respawn watchdog dependency chain

**Build script:** `scripts/build-image.sh <chapter>` → creates Alpine rootfs → packs to ext4 `.img` → uploads to R2

**Acceptance criteria:** Each chapter's flag is ONLY obtainable by solving the puzzle. Images are <200MB each. Total R2 storage <2GB (well within free 10GB tier). CheerpX boots each image in <5s.

---

## Phase 5 — Competition Layer

- [ ] Weekly Quest system:
  - New quest every Friday 8 PM IST, closes Sunday 8 PM IST
  - Cron job rotates quest at exactly 20:00 IST Friday
  - Fresh container per participant (same image for all)
  - Score = time × (1 + 0.1 × hints_used) × (1 + 0.05 × wrong_attempts)
  - Past quests stay playable (unranked)
  - Email reminder Friday 6 PM IST via Gmail SMTP
- [ ] Daily Challenge system:
  - One challenge per day, rotates at midnight IST
  - +50 XP on completion, streak increment
  - No hints, no time limit
  - Streak warning email at 8 PM if not completed
- [ ] Real-time leaderboard via WebSocket — updates within 5s of flag submission

---

## Phase 6 — Polish & Launch Prep

- [ ] Mobile: terminal collapses to bottom sheet, mission brief above on small screens
- [ ] Profile card PNG export (OG image via headless Chrome / Puppeteer or canvas)
- [ ] Elo badge embed snippet (GitHub README markdown generator)
- [ ] Rate limiting on sandbox start (max 3 concurrent sessions per user)
- [ ] `robots.txt`, meta tags for SEO (shell description, og:image of terminal)
- [ ] Load test: 100 concurrent sandbox sessions — p95 latency <500ms
- [ ] Security audit: sandbox escape attempt, flag bypass, auth bypass, chapter skip bypass
- [ ] Pre-warm container pool (5 idle containers per chapter) to reduce cold-start latency

---

## Backlog (Post-Launch)

- [ ] Cyber Heist campaign
- [ ] Mars Colony campaign
- [ ] Corporate Breach campaign
- [ ] Community campaign authoring tool (story JSON + Docker image + validator)
- [ ] Interview Mode (no story, no hints, timed, real SRE problems)
- [ ] Team Battles (college vs college, club competitions)
- [ ] Firecracker microVM migration (replace Docker sandboxes)
- [ ] Command Line Hero replay (cinematic command history after Ch 11 complete)
- [ ] Skill Tree visualization (`cat skillgraph` renders ASCII command graph)
- [ ] Vim Track, Git Track, Kubernetes Track
- [ ] Analytics dashboard (`cat analytics` in shell)

---

## Agent Notes

- Read `story.md` before generating any chapter content or user-facing copy.
- Flags must be stored **hashed** (SHA-256) in the DB — never plaintext, never in API responses.
- The frontend has **no pages** — it's one `xterm.js` instance. No React Router. No URL changes.
- The Google OAuth popup must use `postMessage` to return the JWT to the shell — do not redirect the main window.
- Gmail App Password goes in `.env` as `GMAIL_APP_PASSWORD` — never hardcode.
- When the player runs `submit`, compare SHA-256 of their input against the stored hash — never send the real flag to the client to compare.
- Container progression gate is enforced in **both** the API (server-side) and the shell UX (client-side message) — never trust the client alone.
- The ASCII India map in `cat missions/antariksha/README` must use real city positions relative to each other — not a random layout.
- All user-facing copy must preserve the CTF classified-document tone — no friendly onboarding language.
- `canAccessChapter(userID, chapterNum)` must be called before every `POST /api/sandbox/start` — no exceptions.
