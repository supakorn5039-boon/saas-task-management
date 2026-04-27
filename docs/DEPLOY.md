# Deploying for free

Recommended free-tier stack:

| Layer    | Service        | Free tier                                          |
| -------- | -------------- | -------------------------------------------------- |
| Frontend | **Vercel**     | unlimited bandwidth, free TLS, instant Git deploys |
| Backend  | **Render**     | 750 hrs/mo web service (sleeps after 15 min idle)  |
| Database | **Neon**       | 3 GB Postgres, no sleep, no expiration             |

The cold-start delay on Render (~30s after the backend has been idle) is
the only real catch. For an interview demo, hit your URL once 30 seconds
ahead of time to wake it up.

This guide walks you through all three. **You'll do the clicks**; the
project is already configured to read everything from environment variables
so there's nothing to change in code.

---

## 0. Prerequisites

- Code pushed to GitHub (already done — `origin/main` is current).
- A GitHub account (you'll sign into the three services with it).

---

## 1. Database — Neon (3 minutes)

1. Sign up at https://neon.tech with GitHub.
2. Click **Create project**. Region: pick the one nearest your users.
3. Once created, you'll see a **Connection string** like:
   ```
   postgresql://user:password@ep-xxx.aws.neon.tech/neondb?sslmode=require
   ```
4. **Copy this string** — you'll paste it into Render in the next step.
5. (Optional) Click **SQL Editor** to confirm you can connect. Try
   `SELECT version();`.

That's it. Neon stays free as long as your DB is under 3 GB (it will be).

---

## 2. Backend — Render (5 minutes)

1. Sign up at https://render.com with GitHub.
2. Click **New +** → **Web Service**.
3. Connect your GitHub account. Select the `saas-task-management` repo.
4. Fill in:
   - **Name**: `saas-task-management-api` (or whatever)
   - **Region**: same as your Neon DB
   - **Branch**: `main`
   - **Root Directory**: `backend`
   - **Runtime**: `Docker` (Render auto-detects `dockerfile/build_main.Dockerfile`)
   - **Dockerfile Path**: `dockerfile/build_main.Dockerfile`
   - **Instance Type**: `Free`
5. **Environment variables** — click **Advanced** → add these:

   | Key           | Value                                                                                                       |
   | ------------- | ----------------------------------------------------------------------------------------------------------- |
   | `JWT_SECRET`  | Generate with `openssl rand -hex 32` and paste the result. **Don't commit this anywhere.**                  |
   | `DATABASE_URL`| Paste the Neon connection string from step 1.4.                                                             |
   | `PRODUCTION`  | `true`                                                                                                      |
   | `FRONTEND_URL`| Leave blank for now — you'll come back and set it after Vercel gives you a URL in step 3.                   |

   Render will inject `PORT` automatically — don't set it manually.

6. Click **Create Web Service**. First build takes 5–10 minutes (it builds
   the Docker image and deploys).
7. Once it says "Live" at the top, copy the URL — something like
   `https://saas-task-management-api.onrender.com`.
8. Test: `curl https://saas-task-management-api.onrender.com/api/healthz`
   should return `{"status":"ok","db":"ok"}`.

   (First request takes ~30s as the container wakes; subsequent requests
   are instant within the 15-min idle window.)

---

## 3. Frontend — Vercel (3 minutes)

1. Sign up at https://vercel.com with GitHub.
2. Click **Add New** → **Project**.
3. Import the `saas-task-management` repo.
4. Configure:
   - **Framework Preset**: `Vite` (auto-detected)
   - **Root Directory**: `frontend`
   - **Build Command**: `npm run build` (default, leave alone)
   - **Output Directory**: `dist` (default)
5. **Environment variables** — add one:

   | Key             | Value                                                |
   | --------------- | ---------------------------------------------------- |
   | `VITE_API_URL`  | `https://saas-task-management-api.onrender.com/api`  |

   Use the URL from step 2.7 with `/api` appended.

6. Click **Deploy**. Builds in ~1 minute.
7. Vercel hands you a URL like `https://saas-task-management-xyz.vercel.app`.

---

## 4. Wire CORS so the frontend can call the backend

Now go **back** to Render → your service → **Environment** tab → add:

| Key            | Value                                            |
| -------------- | ------------------------------------------------ |
| `FRONTEND_URL` | Your Vercel URL from step 3.7 (no trailing slash) |

Save. Render will redeploy automatically (~30s).

---

## 5. First-time setup

Open your Vercel URL. You'll get the landing page.

Click **Register**, create yourself an account. That's your admin candidate
— but it'll be created with role `user`. To make it admin, run this in
Neon's SQL Editor (one time):

```sql
UPDATE users SET role = 'admin' WHERE email = 'YOUR_EMAIL_HERE';
```

Refresh, log in. The **Users** menu item now appears in the sidebar. You
have a fully working admin dashboard hosted live, free.

---

## 6. After the first deploy

- **Subsequent pushes to `main`** redeploy both Render and Vercel
  automatically.
- **First request after idle** to the backend takes ~30s. To avoid that
  during a demo, hit `/api/ping` 30 seconds beforehand:
  ```bash
  curl https://saas-task-management-api.onrender.com/api/ping
  ```
- **Logs**: Render's "Logs" tab shows your structured `slog` output. Search
  for `request_id=` to follow a single request through the pipeline.

---

## Troubleshooting

| Symptom                                              | Fix                                                                                                                  |
| ---------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Backend logs `jwt_secret must be at least 32 chars`  | Your `JWT_SECRET` env var is too short. Regenerate: `openssl rand -hex 32` is 64 chars and works.                    |
| `/api/healthz` returns 503 with `db: ping_failed`    | `DATABASE_URL` is wrong, missing, or the Neon DB is paused. Copy the string fresh from Neon — note `?sslmode=require`.|
| Frontend gets `CORS error` calling backend           | `FRONTEND_URL` on Render doesn't match your Vercel URL exactly. No trailing slash. Save → wait for redeploy.         |
| Frontend gets 404 on `/api/...`                      | `VITE_API_URL` on Vercel is wrong or missing the `/api` suffix.                                                      |
| First page load is slow                              | Backend cold-starting after idle. Normal. ~30s. Subsequent requests are fast.                                        |

---

## Tearing it down

- Delete the Render service → free tier resets immediately.
- Delete the Neon project → 3 GB freed.
- Delete the Vercel project → instant.

Nothing here charges you a cent at any point on the free tier — but verify
your billing settings on each service have no card attached if you want to
be 100% sure.
