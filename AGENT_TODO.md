# LinuxQuest — Agent TODO

> **For AI Agents:** This file is your source of truth. Read it fully before taking any action.
> Mark tasks `[/]` when starting, `[x]` when done. Never skip acceptance criteria.
> Always check `## Context` and `## Constraints` before writing code.

---

## Context

**Project:** LinuxQuest — a CTF-style, story-driven Linux learning platform.  
**Flagship campaign:** Operation Antariksha (story in `story.md`, overview in `README.md`).  
**Stack (planned):**
- Frontend: React + Vite (TypeScript)
- Backend: Go + Gin
- Database: PostgreSQL
- Container runtime: Docker (Firecracker later)
- Auth: JWT
- Realtime: WebSockets (leaderboard, live terminal)

**Key files to read before coding:**
- `README.md` — platform overview and mission map
- `story.md` — full chapter dossiers, flag formats, narrative
- `orignalreadme.md` — original architecture notes and file split suggestions

---

## Constraints

- Do NOT skip phases — each phase depends on the previous.
- Do NOT use placeholder data in the UI — generate real content from `story.md`.
- Every new feature needs a working test before marking `[x]`.
- Keep the CTF mystery tone in all user-facing copy.
- All containers must be sandboxed — no host filesystem access.

---

## Phase 1 — Foundation & Repo Setup

### 1.1 Project Scaffold
- [ ] Initialize monorepo structure:
  ```
  linuxquest/
  ├── app/          ← React + Vite frontend
  ├── server/       ← Go + Gin backend
  ├── sandbox/      ← Docker sandbox templates per chapter
  ├── campaigns/    ← JSON campaign definitions
  ├── docs/         ← Architecture, API docs
  └── scripts/      ← Dev utilities
  ```
- [ ] Add root `Makefile` with targets: `dev`, `build`, `test`, `lint`, `sandbox`
- [ ] Add `.env.example` with all required env vars documented
- [ ] Add `docker-compose.yml` for local dev (postgres, backend, frontend, sandbox)
- [ ] Add GitHub Actions CI: lint + test on every PR

**Acceptance criteria:** `make dev` starts the full stack locally. `make test` passes with 0 failures.

---

### 1.2 Database Schema
- [ ] Design and migrate initial schema:
  - `users` — id, username, email, password_hash, elo, level, xp, streak, created_at
  - `campaigns` — id, slug, title, description, difficulty
  - `chapters` — id, campaign_id, number, title, city, difficulty, commands[], story_text, flag_hash
  - `user_progress` — user_id, chapter_id, status (locked/active/complete), attempts, completed_at
  - `submissions` — id, user_id, chapter_id, flag_submitted, is_correct, submitted_at
  - `sessions` — id, user_id, container_id, started_at, expires_at
  - `leaderboard` — weekly snapshot table
- [ ] Write seed data for Operation Antariksha (Ch 0–11) from `story.md`
- [ ] Write migration scripts (up + down)

**Acceptance criteria:** `make db-migrate` runs clean. Seed data loads all 12 chapters with correct flags.

---

## Phase 2 — Backend API

### 2.1 Auth
- [ ] `POST /api/auth/register` — create user, hash password (bcrypt), return JWT
- [ ] `POST /api/auth/login` — validate credentials, return JWT + refresh token
- [ ] `POST /api/auth/refresh` — refresh access token
- [ ] JWT middleware — attach user to context on all protected routes
- [ ] Rate limiting on auth endpoints (max 5 req/min per IP)

**Acceptance criteria:** Register → login → access protected route works end-to-end. Invalid tokens return 401.

---

### 2.2 Campaign & Chapter API
- [ ] `GET /api/campaigns` — list all campaigns with metadata
- [ ] `GET /api/campaigns/:slug` — campaign detail + chapter list
- [ ] `GET /api/campaigns/:slug/chapters/:num` — chapter detail (story text, commands, difficulty) — **no flag in response**
- [ ] `POST /api/campaigns/:slug/chapters/:num/submit` — validate flag submission (compare hash), update progress, award XP
- [ ] `GET /api/users/:username/progress` — return all chapter statuses for a user

**Acceptance criteria:** Flag submission endpoint correctly validates correct flags, rejects wrong ones, and is not bypassable by submitting the hash directly.

---

### 2.3 Sandbox API
- [ ] `POST /api/sandbox/start` — spin up a Docker container for a chapter, return `session_id` + WebSocket URL
- [ ] `DELETE /api/sandbox/:session_id` — stop and remove container
- [ ] `GET /api/sandbox/:session_id/status` — check container health
- [ ] Auto-expire containers after 45 minutes of inactivity
- [ ] Each container must be:
  - Network-isolated (no outbound internet)
  - Read-only root filesystem except `/tmp` and `/home/player`
  - CPU-limited (0.5 cores), Memory-limited (256 MB)
  - Pre-loaded with chapter-specific files from `sandbox/<chapter>/` directory

**Acceptance criteria:** Starting a sandbox returns a working WebSocket terminal. Container cannot reach the internet. Container auto-destroys after timeout.

---

### 2.4 Leaderboard & Elo API
- [ ] `GET /api/leaderboard/weekly` — top 50 by weekly XP
- [ ] `GET /api/leaderboard/global` — top 50 by Elo
- [ ] `GET /api/users/:username/elo` — return Elo + rank band
- [ ] Elo update logic: beating a Hard chapter = larger Elo gain; hints used = penalty
- [ ] Elo bands: Newcomer (<800), User (800–1199), Power User (1200–1599), Sysadmin (1600–1999), Engineer (2000–2399), Wizard (2400+)

**Acceptance criteria:** Completing a Hard chapter increases Elo by more than Standard. Using a hint reduces the Elo gain.

---

## Phase 3 — Frontend

### 3.1 Design System
- [ ] Set up global CSS variables:
  - Colors: dark background (#0d1117), accent green (#39ff14), red (#ff4444), amber (#ffaa00)
  - Font: `JetBrains Mono` for terminal UI, `Inter` for prose
  - Spacing, radius, shadow tokens
- [ ] Component library (Storybook optional):
  - `<Terminal />` — xterm.js based, WebSocket-connected
  - `<MissionBrief />` — classified dossier card with blinking cursor animation
  - `<ChapterMap />` — India map with 11 city nodes, color-coded by status
  - `<EloCard />` — user rating with rank band badge
  - `<ProgressBar />` — animated XP bar
  - `<FlagInput />` — monospace input with `ANTARIKSHA{}` prefix locked

**Acceptance criteria:** Each component renders correctly in isolation with mock data.

---

### 3.2 Pages
- [ ] `/` — Landing page: CTF-style briefing, SHIVA intro, "Start Mission" CTA
- [ ] `/login` and `/register` — minimal, dark-themed auth forms
- [ ] `/campaign/antariksha` — campaign overview with India chapter map
- [ ] `/campaign/antariksha/chapter/:num` — chapter page:
  - Left: `<MissionBrief />` with story text, commands, and flag input
  - Right: `<Terminal />` connected to sandbox WebSocket
  - Top bar: timer, hints remaining, chapter title + city
- [ ] `/profile/:username` — Elo card, skill breakdown bars, badges, progress grid
- [ ] `/leaderboard` — weekly and global tabs, live-updating via WebSocket

**Acceptance criteria:** Full chapter flow works end-to-end: land → register → start chapter → terminal opens → submit flag → XP awarded → progress updated.

---

### 3.3 Terminal Integration
- [ ] Integrate `xterm.js` with `xterm-addon-fit` and `xterm-addon-web-links`
- [ ] WebSocket connection to backend sandbox endpoint
- [ ] On connect: display chapter-specific welcome message (from `story.md` mission brief)
- [ ] On disconnect: show "Connection lost — reconnecting…" with retry logic
- [ ] Copy-paste support, proper ctrl+c handling

**Acceptance criteria:** Player can type and run real Linux commands in the browser terminal against the sandboxed container.

---

## Phase 4 — Sandbox Content (Chapter Files)

For each chapter (0–11), create `sandbox/chapter-<N>/` with:

- [ ] **Ch 0:** Basic filesystem with breadcrumbs hidden in subdirectories (`hint_*.txt` files)
- [ ] **Ch 1:** Scattered files, mislabeled directories, partially overwritten incident log with SHIVA payload hidden inside
- [ ] **Ch 2:** 500,000-line telemetry log with `SHIVA_PING` embedded in exactly 60 lines
- [ ] **Ch 3:** Three running processes — one SHIVA, two legitimate firewall daemons with similar names
- [ ] **Ch 4:** Malicious cron file in `/etc/cron.d/` mimicking a legitimate package name
- [ ] **Ch 5:** SUID binary modified to grant root shell to any user
- [ ] **Ch 6:** Staged `.tar.gz` archive + baseline backup for diff
- [ ] **Ch 7:** 2.3M-line mixed JSON/plaintext proxy log with 47 matching POST requests
- [ ] **Ch 8:** 14-node cluster simulation (Docker network) with shared NFS mount worm loop
- [ ] **Ch 9:** Live DNS traffic with anomalous query pattern on `eth1`
- [ ] **Ch 10:** 3-hop SSH jump chain to air-gapped target node
- [ ] **Ch 11:** Systemd service unit with respawn watchdog and dependency chain

**Acceptance criteria:** Each chapter's flag is only obtainable by solving the puzzle — not by guessing or reading source files directly.

---

## Phase 5 — Competition Layer

- [ ] Weekly Quest system:
  - New quest every Friday 8 PM IST → Sunday 8 PM IST
  - Fresh container for each participant
  - Score = time × hint penalty × wrong-attempt penalty
  - Past quests stay playable (unranked)
- [ ] Daily Challenge system:
  - One challenge per day, auto-rotated at midnight IST
  - +50 XP on completion, streak increment
  - No hints
- [ ] Real-time leaderboard via WebSocket (updates on flag submission)

---

## Phase 6 — Polish & Launch Prep

- [ ] Mobile responsiveness (terminal collapses to bottom sheet on small screens)
- [ ] Profile card PNG export (OG image generation via canvas or Puppeteer)
- [ ] Elo badge embed code (GitHub README snippet generator)
- [ ] Rate limiting and abuse protection on sandbox start endpoint
- [ ] Error boundary on terminal WebSocket disconnect
- [ ] `robots.txt`, `sitemap.xml`, meta tags for SEO
- [ ] Load test: 100 concurrent sandbox sessions
- [ ] Security audit: sandbox escape attempt, flag bypass attempt, auth bypass attempt

---

## Backlog (Post-Launch)

- [ ] Cyber Heist campaign
- [ ] Mars Colony campaign
- [ ] Corporate Breach campaign
- [ ] Community campaign authoring tool
- [ ] Interview Mode (no story, no hints, timed)
- [ ] Team Battles
- [ ] Firecracker microVM migration (replace Docker sandboxes)
- [ ] Vim Track, Git Track, Kubernetes Track

---

## Agent Notes

- Read `story.md` before generating any chapter content or copy.
- The flag for each chapter must be stored **hashed** (SHA-256) in the DB — never plaintext.
- When spinning up sandbox containers, pull from pre-built images in `sandbox/<chapter>/Dockerfile` — do not build at runtime.
- The India chapter map (`<ChapterMap />`) should use actual city coordinates (lat/lng), not a static image.
- All user-facing text must preserve the CTF classified-document tone — no friendly onboarding language.
