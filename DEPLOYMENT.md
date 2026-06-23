# LinuxQuest — Deployment Guide

> **Total cost: ₹0/month** (all free tiers)
> Time to deploy: ~45 minutes for first setup.

---

## Architecture Recap

```
Browser (player)
├── xterm.js shell UI          → Vercel (free)
├── CheerpX WASM sandbox       → Cloudflare R2 disk images (free)
└── API calls (auth, progress) → Go backend on Render.com (free)
                                      ↓
                               Supabase PostgreSQL (free)
                               Gmail SMTP (free)
                               Google OAuth (free)
```

---

## Prerequisites

Install these locally before starting:

```bash
# Node.js 20+
node --version

# Go 1.22+
go version

# Wrangler (Cloudflare CLI)
npm install -g wrangler
wrangler login
```

---

## Step 1 — Google OAuth Setup

1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Create a new project → **LinuxQuest**
3. APIs & Services → OAuth consent screen → External → Fill in app name
4. APIs & Services → Credentials → Create OAuth 2.0 Client ID
   - Application type: **Web application**
   - Authorized redirect URIs:
     - `http://localhost:8080/api/auth/google/callback` (local dev)
     - `https://linuxquest-api.onrender.com/api/auth/google/callback` (production)
5. Copy **Client ID** and **Client Secret** → save for Step 4

---

## Step 2 — Gmail App Password Setup

1. Go to [myaccount.google.com/security](https://myaccount.google.com/security)
2. Enable **2-Step Verification** (required)
3. Search for **App passwords**
4. Generate → App name: `LinuxQuest`
5. Copy the 16-character password → save for Step 4

> ⚠️ This is different from your Gmail password. Keep it secret.

---

## Step 3 — Supabase (PostgreSQL)

1. Go to [supabase.com](https://supabase.com) → New project → **linuxquest**
2. Choose a strong database password → save it
3. Wait ~2 minutes for provisioning
4. Go to **Project Settings → Database**
5. Copy the **Connection string (URI)** → save for Step 4

Run migrations after deploy:
```bash
# From project root
go run ./scripts/migrate up
```

---

## Step 4 — Environment Variables

Create `.env` in the project root (never commit this):

```env
# Google OAuth
GOOGLE_CLIENT_ID=your_client_id_here
GOOGLE_CLIENT_SECRET=your_client_secret_here
GOOGLE_REDIRECT_URI=https://linuxquest-api.onrender.com/api/auth/google/callback

# Gmail SMTP
GMAIL_USER=your_email@gmail.com
GMAIL_APP_PASSWORD=xxxx xxxx xxxx xxxx

# Database
DATABASE_URL=postgresql://postgres:password@db.xxx.supabase.co:5432/postgres

# JWT
JWT_SECRET=generate_a_random_64_char_string_here

# Cloudflare R2
R2_ACCOUNT_ID=your_account_id
R2_ACCESS_KEY_ID=your_access_key
R2_SECRET_ACCESS_KEY=your_secret_key
R2_BUCKET_NAME=linuxquest-images
R2_PUBLIC_URL=https://pub-xxx.r2.dev

# App
PORT=8080
ENV=production
FRONTEND_URL=https://linuxquest.vercel.app
```

Generate a secure JWT secret:
```bash
openssl rand -hex 32
```

---

## Step 5 — Cloudflare R2 (Disk Images)

1. Go to [dash.cloudflare.com](https://dash.cloudflare.com) → R2 → Create bucket
   - Bucket name: `linuxquest-images`
   - Region: Auto
2. Go to R2 → Manage R2 API Tokens → Create API Token
   - Permissions: Object Read & Write
   - Copy Access Key ID and Secret Access Key → add to `.env`
3. Enable public access: R2 → `linuxquest-images` → Settings → Public Access → Allow

Build and upload disk images (after writing chapter content in `images/`):
```bash
# Build one chapter
make image chapter=1

# Build and upload all chapters
make images-all

# Or manually
./scripts/build-image.sh 1   # builds images/ch1.img
wrangler r2 object put linuxquest-images/ch1.img --file images/ch1.img
```

---

## Step 6 — Backend Deployment (Render.com)

Render provides a completely card-free free tier for Go applications.

1. Create a free account at [render.com](https://render.com) using your GitHub account.
2. From the Dashboard, click **New +** → **Web Service**.
3. Connect your GitHub repository containing the `linuxquest` code.
4. Fill in the Web Service configuration:
   - **Name**: `linuxquest-api`
   - **Region**: Select a close region (e.g., `Singapore` or `Oregon`)
   - **Branch**: `main`
   - **Runtime**: `Go`
   - **Build Command**: `cd server && go build -o main main.go`
   - **Start Command**: `./server/main`
   - **Instance Type**: Select **Free** ($0/month - no credit card required)
5. Expand the **Advanced** section to add the following **Environment Variables**:
   - `GOOGLE_CLIENT_ID` = `your_google_client_id`
   - `GOOGLE_CLIENT_SECRET` = `your_google_client_secret`
   - `GOOGLE_REDIRECT_URI` = `https://linuxquest-api.onrender.com/api/auth/google/callback`
   - `GMAIL_USER` = `your_email@gmail.com`
   - `GMAIL_APP_PASSWORD` = `your_gmail_app_password`
   - `DATABASE_URL` = `your_supabase_postgresql_connection_string`
   - `JWT_SECRET` = `your_random_64_character_jwt_secret`
   - `R2_ACCOUNT_ID` = `your_cloudflare_r2_account_id`
   - `R2_ACCESS_KEY_ID` = `your_r2_access_key_id`
   - `R2_SECRET_ACCESS_KEY` = `your_r2_secret_key`
   - `R2_BUCKET_NAME` = `linuxquest-images`
   - `R2_PUBLIC_URL` = `https://pub-xxx.r2.dev`
   - `FRONTEND_URL` = `https://linuxquest.vercel.app`
   - `ENV` = `production`
6. Click **Create Web Service**.

Render will automatically fetch, build, and deploy your Go backend. 
Your API will be live at:
`https://linuxquest-api.onrender.com`

> [!NOTE]
> Since this is on Render's Free tier, the server will automatically go to sleep after 15 minutes of inactivity. When a player visits the app after it goes to sleep, the first request might take 50-60 seconds to wake the service back up.

To run migrations, you can execute a one-off migration script from your local machine targeting the Supabase production DB connection string, or run it inside the Go server boot sequence.


---

## Step 7 — Frontend Deployment (Vercel)

```bash
# From the app/ directory
cd app

# Install Vercel CLI
npm install -g vercel

# Deploy (first time — follow prompts)
vercel

# Set environment variables on Vercel
vercel env add VITE_API_URL
# Enter: https://linuxquest-api.onrender.com

vercel env add VITE_GOOGLE_CLIENT_ID
# Enter: your_google_client_id

vercel env add VITE_R2_PUBLIC_URL
# Enter: https://pub-xxx.r2.dev

# Deploy to production
vercel --prod

# Your frontend is live at:
# https://linuxquest.vercel.app
```

Or connect GitHub repo to Vercel for automatic deploys on every push:
1. [vercel.com](https://vercel.com) → New Project → Import GitHub repo
2. Root Directory: `app/`
3. Build Command: `npm run build`
4. Output Directory: `dist`
5. Add env vars in Vercel dashboard → Settings → Environment Variables

---

## Step 8 — Update Google OAuth Redirect URI

Go back to [console.cloud.google.com](https://console.cloud.google.com):
- APIs & Services → Credentials → your OAuth client
- Add to Authorized redirect URIs:
  - `https://linuxquest-api.onrender.com/api/auth/google/callback`
- Add to Authorized JavaScript origins:
  - `https://linuxquest.vercel.app`

---

## Local Development

```bash
# Clone repo
git clone https://github.com/divyansh-v15-06/linux
cd linuxquest

# Copy env
cp .env.example .env
# Fill in .env with your values

# Start everything
make dev
# Frontend: http://localhost:5173
# Backend:  http://localhost:8080

# Or manually:
# Terminal 1 — Backend
cd server && go run .

# Terminal 2 — Frontend
cd app && npm install && npm run dev
```

For local OAuth, make sure `http://localhost:8080/api/auth/google/callback` is in your Google OAuth redirect URIs.

---

## Verifying the Deployment

After all steps, check:

```bash
# 1. Backend health
curl https://linuxquest-api.onrender.com/api/health
# Expected: {"status":"ok"}

# 2. R2 image accessible
curl -I https://pub-xxx.r2.dev/ch0.img
# Expected: HTTP/2 200

# 3. Frontend
# Open https://linuxquest.vercel.app
# Expected: ISRO-CIRT terminal boot sequence
```

Full flow test:
```
1. Open linuxquest.vercel.app
2. Boot sequence plays → prompt appears
3. Type: login → Google OAuth popup opens
4. Authenticate → prompt changes to arjun@linuxquest
5. Type: cd missions/antariksha/ch0
6. Type: start → CheerpX boots Alpine Linux
7. Type: ls → files appear
8. Type: exit → returns to mission shell
```

---

## Continuous Deployment

Push to `main` → auto-deploys to both Vercel and Fly.io:

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy --remote-only
        working-directory: server/
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}

  deploy-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: npm ci && npm run build
        working-directory: app/
      - uses: amondnet/vercel-action@v25
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.VERCEL_ORG_ID }}
          vercel-project-id: ${{ secrets.VERCEL_PROJECT_ID }}
          working-directory: app/
```

Add secrets in GitHub → Settings → Secrets:
- `FLY_API_TOKEN` → from `flyctl auth token`
- `VERCEL_TOKEN` → from Vercel dashboard → Settings → Tokens
- `VERCEL_ORG_ID` → from `vercel whoami`
- `VERCEL_PROJECT_ID` → from `.vercel/project.json` after first deploy

---

## Cost Summary

| Service | Free Tier Limits | Your Usage |
|---------|-----------------|------------|
| **Vercel** | 100GB bandwidth/month | ~1–5GB (static files) ✅ |
| **Fly.io** | 3 shared VMs, 256MB RAM | 1 VM, 256MB ✅ |
| **Supabase** | 500MB DB, 5GB bandwidth | ~50MB (users + progress) ✅ |
| **Cloudflare R2** | 10GB storage, 10M reads | ~2GB (12 images) ✅ |
| **Google OAuth** | Unlimited | Unlimited ✅ |
| **Gmail SMTP** | 500 emails/day | ~50/day (beta) ✅ |
| **Total** | | **₹0/month** |

---

## Troubleshooting

| Problem | Fix |
|---------|-----|
| OAuth redirect mismatch | Check redirect URI in Google Cloud Console matches exactly |
| Gmail emails going to spam | Set up SPF record in your domain DNS (or use a real domain later) |
| Fly.io app sleeping | Free tier auto-stops after inactivity — first request takes ~5s to wake |
| CheerpX not loading | Check R2 CORS settings — allow `GET` from `https://linuxquest.vercel.app` |
| DB migration fails | Check `DATABASE_URL` is correct; Supabase requires SSL (`?sslmode=require`) |
| R2 image 403 | Check bucket is set to public and R2 token has read permissions |
