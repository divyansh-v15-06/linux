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
- Backend: Go + Gin
- Database: PostgreSQL (Supabase free tier)
- Auth: Google OAuth 2.0 + JWT
- Email: Gmail SMTP with App Password
- Sandbox: Docker → Firecracker (future)
- Realtime: WebSockets (terminal I/O, leaderboard)

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
- Do NOT build at container runtime — pre-build Docker images in `sandbox/<chapter>/`.
- Every new feature needs a working test before marking `[x]`.
- All containers must be sandboxed — no host filesystem access, no outbound internet.

---

## Phase 1 — Foundation & Repo Setup

### 1.1 Project Scaffold
- [ ] Initialize monorepo structure:
  ```
  linuxquest/
  ├── app/          ← React + Vite frontend (shell UI)
  ├── server/       ← Go + Gin backend
  ├── sandbox/      ← Docker images per chapter
  ├── campaigns/    ← JSON campaign + chapter definitions
  ├── docs/         ← Architecture, API docs
  └── scripts/      ← Dev utilities, seed scripts
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

### 2.4 Sandbox API
- [ ] `POST /api/sandbox/start` — check progression gate, spin up Docker container for chapter, return `session_id` + WebSocket URL
- [ ] `DELETE /api/sandbox/:session_id` — stop and remove container
- [ ] `GET /api/sandbox/:session_id/status` — check container health
- [ ] Auto-expire containers after 30 minutes of inactivity
- [ ] Container constraints (enforced at `docker run`):
  - `--network none` (Ch 1–8, 11) or isolated internal net (Ch 9–10)
  - `--read-only` root FS with `--tmpfs /tmp --tmpfs /home/player`
  - `--cpus 0.5 --memory 128m`
  - `--security-opt=no-new-privileges --cap-drop ALL`
- [ ] Pre-pull chapter images at server startup — never build at runtime

**Acceptance criteria:** Sandbox starts in <3s. Container cannot reach the internet. Auto-destroys after timeout. `rm -rf /` inside container is harmless — new container spawns for next attempt.

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

**Acceptance criteria:** Site loads as full-screen black terminal. Boot sequence plays. Prompt appears. No other UI visible.

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
- [ ] `start` calls `POST /api/sandbox/start`, opens WebSocket, switches terminal to live container I/O
- [ ] `exit` / `Ctrl+D` inside sandbox closes WebSocket, returns to mission shell
- [ ] `submit ANTARIKSHA{...}` calls submit API, prints green (correct) or red (wrong) result

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

### 3.6 WebSocket Sandbox Integration
- [ ] Connect to `ws://backend/sandbox/:session_id` when player runs `start`
- [ ] Pipe all xterm.js stdin → WebSocket → container stdin
- [ ] Pipe all container stdout → WebSocket → xterm.js stdout
- [ ] On disconnect: print `[CONNECTION LOST — reconnecting...]` in amber, auto-retry 3×
- [ ] On container timeout: print `[SESSION EXPIRED — type  cd ch<N>  to restart]`
- [ ] `Ctrl+D` / `exit` gracefully closes WebSocket and returns to mission shell

**Acceptance criteria:** Player types real Linux commands in browser, they execute in Docker container, output appears in <200ms.

---

## Phase 4 — Sandbox Content (Chapter Files)

For each chapter, create `sandbox/chapter-<N>/` with a `Dockerfile` and pre-seeded filesystem:

- [ ] **Ch 0:** Basic filesystem with `hint_*.txt` breadcrumbs hidden in subdirectories
- [ ] **Ch 1:** Scattered files, mislabeled directories, partially overwritten incident log with SHIVA payload
- [ ] **Ch 2:** 500,000-line telemetry log with `SHIVA_PING` in exactly 60 lines
- [ ] **Ch 3:** Three running processes — one SHIVA (high CPU), two legitimate firewall daemons with similar names
- [ ] **Ch 4:** Malicious cron file in `/etc/cron.d/` mimicking a legitimate package name
- [ ] **Ch 5:** Modified SUID binary that grants root shell to any user
- [ ] **Ch 6:** Staged `.tar.gz` archive + `/backup/baseline.tar.gz` for diff
- [ ] **Ch 7:** 2.3M-line mixed JSON/plaintext proxy log with 47 matching POST requests with steganographic User-Agent
- [ ] **Ch 8:** 14-node Docker network simulation with shared NFS mount worm loop
- [ ] **Ch 9:** `tcpdump` on `eth1` with live anomalous DNS query pattern to C2 server
- [ ] **Ch 10:** 3-hop SSH jump chain (`jump1 → jump2 → target`) with different keys per hop
- [ ] **Ch 11:** systemd service with respawn watchdog + dependency chain — must stop in correct order

**Acceptance criteria:** Each chapter's flag is ONLY obtainable by solving the puzzle. Direct file reads cannot reveal the flag. All images build in <60s. All images <500MB.

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
