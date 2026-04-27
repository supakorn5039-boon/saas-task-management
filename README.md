# SaaS Task Management

A full-stack task tracker built end-to-end as a portfolio piece — Go on the
back, React 19 on the front, Postgres in the middle. Built to demonstrate
production patterns I'd actually ship at work, not a bootcamp todo app.

[![CI](https://github.com/supakorn5039-boon/saas-task-management/actions/workflows/ci.yml/badge.svg)](https://github.com/supakorn5039-boon/saas-task-management/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?logo=typescript&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.5-4169E1?logo=postgresql&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-4-06B6D4?logo=tailwindcss&logoColor=white)
![Tests](https://img.shields.io/badge/Tests-passing-brightgreen)

---

## What this project actually is

A multi-user task management app where:

- **Users** sign up, log in, manage their own tasks (CRUD with full edit, status flips, search, filter, sort, paginated table).
- **Admins** manage other users — change role, deactivate, delete (with self-protection so an admin can't lock themselves out).
- **Account settings** — change own password with current-password verification.
- **Dashboard** shows live stats and recent activity from a single API call.

What makes it more than a tutorial:

- **Versioned, reversible SQL migrations** (not GORM AutoMigrate) — see [ADR 001](docs/decisions/001-raw-sql-migrations.md).
- **Structured `AppError` envelope** — services classify errors, controllers map cleanly to status codes; GORM strings never leak to clients. See [ADR 002](docs/decisions/002-apperror-envelope.md).
- **TanStack Query key factories** — single source of truth for cache invalidation. See [ADR 003](docs/decisions/003-query-key-factories.md).
- **Optimistic UI updates** on the most-common interaction (task status flip) with full rollback on failure.
- **Real graceful shutdown** — SIGTERM drains in-flight requests, doesn't drop them.
- **Real `/healthz`** — pings the DB, returns 503 if Postgres is unreachable.
- **Per-IP rate limiting** on the auth endpoints (in-memory token-bucket).
- **Conservative security headers** — XCTO, X-Frame-Options, Referrer-Policy.
- **Structured request logging** with `slog` and per-request IDs propagated via `X-Request-ID`.
- **JWT secret enforcement** — fails fast at boot if missing or under 32 chars; supports `JWT_SECRET` env override.
- **Tests on both sides** — Go service tests against a real Postgres test DB, Vitest + React Testing Library on the frontend.
- **CI on every push** via GitHub Actions — lint, type-check, test, build for both stacks.

---

## Tech stack

| Layer        | Tech                                                              |
| ------------ | ----------------------------------------------------------------- |
| **Frontend** | React 19, TypeScript 5.9, Vite 8, TanStack Router/Query/Table, Zustand, react-hook-form, Zod |
| **Backend**  | Go 1.25, Gin, GORM, JWT (golang-jwt v5), bcrypt                   |
| **Database** | PostgreSQL 17.5 (Docker for local dev)                            |
| **Styling**  | Tailwind CSS 4, shadcn/ui, lucide-react, dark mode via next-themes |
| **Testing**  | `go test` + Postgres test DB, Vitest + Testing Library + jsdom    |
| **Quality**  | ESLint 9, Prettier, lefthook (pre-commit/pre-push), GitHub Actions CI |
| **Tooling**  | Docker Compose, Bruno (API collection), Vite, slog                |

---

## API

### Public

| Method | Endpoint              | Notes                                  |
| ------ | --------------------- | -------------------------------------- |
| GET    | `/api/ping`           | Liveness                               |
| GET    | `/api/healthz`        | Readiness — pings DB                   |
| POST   | `/api/auth/login`     | Rate-limited (5/min/IP)                |
| POST   | `/api/auth/register`  | Rate-limited; password ≥ 8 chars       |

### Authenticated

| Method | Endpoint               | Notes                                                     |
| ------ | ---------------------- | --------------------------------------------------------- |
| GET    | `/api/user/profile`    | Current user                                              |
| PUT    | `/api/user/password`   | Change own password (verifies current)                    |
| GET    | `/api/tasks`           | Paginated list, filter, search, sort + status counts      |
| POST   | `/api/tasks`           | Create                                                    |
| PUT    | `/api/tasks/:id`       | Patch any subset of `{title, description, status}`        |
| DELETE | `/api/tasks/:id`       | Soft delete                                               |

### Admin (RBAC enforced)

| Method | Endpoint                   | Notes                                                |
| ------ | -------------------------- | ---------------------------------------------------- |
| GET    | `/api/admin/users`         | Paginated list with search + sort (admin-only)       |
| PUT    | `/api/admin/users/:id`     | Update `role` and/or `status` (self-demote blocked)  |
| DELETE | `/api/admin/users/:id`     | Soft delete (self-delete blocked)                    |

#### `GET /api/tasks` query params

| Param      | Default      | Notes                                                            |
| ---------- | ------------ | ---------------------------------------------------------------- |
| `page`     | `1`          | ≥ 1                                                              |
| `per_page` | `10`         | Clamped 1–100                                                    |
| `status`   | (none)       | `todo` \| `in_progress` \| `done`                                |
| `search`   | (none)       | ILIKE on title + description                                     |
| `sort`     | `created_at` | Whitelisted: `created_at`, `updated_at`, `title`, `status`       |
| `order`    | `desc`       | `asc` \| `desc`                                                  |

Response envelope:

```json
{
  "data":  [{ "id": 1, "title": "...", "status": "todo", "createdAt": "..." }],
  "meta":  { "page": 1, "perPage": 10, "total": 47 },
  "counts":{ "all": 47, "todo": 30, "in_progress": 10, "done": 7 }
}
```

Error envelope (consistent across all endpoints):

```json
{ "error": "task not found" }
```

---

## Project structure

```
saas-task-management/
├── backend/
│   ├── src/
│   │   ├── apiwebserver/
│   │   │   ├── controller/      HTTP handlers (auth, user, task, admin, response helper)
│   │   │   ├── service/         Business logic + GORM queries (+ tests)
│   │   │   └── middleware/      Protected (JWT), Rbac, RequestLogger, RateLimit, SecurityHeaders
│   │   ├── apperror/            Typed error envelope (with tests)
│   │   ├── config/              INI loader + env-var override + JWT secret validation
│   │   ├── database/
│   │   │   ├── database.go      Connect + auto-run migrations
│   │   │   ├── model/           GORM models + DTOs
│   │   │   ├── migration/       Versioned raw-SQL migrations
│   │   │   └── seeder/          Registry-based seeders
│   │   ├── pkg/                 API server bootstrap (router, CORS, healthz)
│   │   ├── security/            JWT + bcrypt
│   │   ├── testhelpers/         Test DB setup
│   │   └── main.go              Entry point + CLI commands + graceful shutdown
│   ├── bruno/                   API collection (Auth, User, Tasks, Admin, Health)
│   ├── docker-compose.yml
│   └── config.example.ini       Copy to config.ini and edit
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── ui/              shadcn primitives
│   │   │   ├── layout/          AppShell, Sidebar, UserMenu
│   │   │   ├── shared/          PageError, Skeletons, StatCard, ConfirmDialog, ThemeToggle, data-table primitives
│   │   │   ├── tasks/           Task domain components
│   │   │   └── admin/           Admin domain components
│   │   ├── routes/              File-based pages (TanStack Router)
│   │   ├── services/            API clients + query-key factories
│   │   ├── hooks/               TanStack Query hooks + list-state helpers
│   │   ├── store/               Zustand auth store
│   │   ├── types/               Cross-feature TS types
│   │   ├── constants/           Display labels
│   │   ├── styles/              Tailwind class maps per domain
│   │   ├── validators/          Zod schemas
│   │   ├── lib/                 axios, api-error, date, router shim, utils
│   │   └── __tests__/           Mirrors src/ — Vitest tests
│   ├── vitest.config.ts
│   ├── package.json
│   └── vite.config.ts
│
├── docs/decisions/              Architecture Decision Records (ADRs)
├── .github/workflows/ci.yml     Lint + type-check + test + build, both stacks
├── lefthook.yml                 Pre-commit / pre-push hook config
└── README.md
```

The conventions for *what goes where* are pinned in [ADR 006](docs/decisions/006-frontend-layer-layout.md).

---

## Getting started

### Prerequisites

- **Go** 1.25+
- **Node** 22+
- **Docker** (for Postgres)
- (Optional) **Lefthook** — `brew install lefthook` to enable git hooks

### 1. Start the database

```bash
cd backend
docker compose up -d
```

### 2. Configure the backend

```bash
cp config.example.ini config.ini
# Edit jwt_secret (≥ 32 chars) and DB credentials.
# Or set JWT_SECRET / DB_PASSWORD as environment variables — they win over config.ini.
```

### 3. Run the backend

```bash
cd backend
go run ./src
```

Backend serves on `http://localhost:8080`. Migrations auto-run on every boot.

CLI commands:

```bash
go run ./src                  # start the HTTP server
go run ./src seed             # seed default roles + users
go run ./src migrate:status   # show migration status
go run ./src migrate:rollback # roll back the last migration
```

### 4. Run the frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend serves on `http://localhost:5173`.

### 5. (Optional) Enable git hooks

```bash
lefthook install
```

Wires up pre-commit (lint, type-check, build per file) and pre-push (full prod build, secret scan, TODO scan).

---

## Default seeded users

Run `cd backend && go run ./src seed` to create:

| Email                | Role  | Password      |
| -------------------- | ----- | ------------- |
| `admin@example.com`  | admin | `password123` |
| `user@example.com`   | user  | `password123` |

Roles control sidebar visibility — `admin` sees the **Users** menu item; `user` does not.

---

## Testing

### Backend

```bash
cd backend
# Make sure docker compose is up + create the test DB once:
PGPASSWORD=password psql -h localhost -U username -d postgres -c "CREATE DATABASE saas_test;"

go test -race -p 1 ./...   # -p 1 prevents parallel package runs from racing on the shared test DB
```

Tests live next to the code (Go convention) and run against a real Postgres
instance via `src/testhelpers/db.go`. The helper auto-migrates the schema and
truncates between tests. If the test DB is unreachable, tests `t.Skip` rather
than fail — so you can still run other test suites without Postgres.

CI uses a Postgres service container (see `.github/workflows/ci.yml`).

### Frontend

```bash
cd frontend
npm test          # one-off run
npm run test:watch
npm run test:ui
```

Tests live in `frontend/src/__tests__/` and mirror the source tree. See
[ADR 004](docs/decisions/004-test-layout.md) for why backend co-locates and
frontend doesn't.

---

## Production-readiness checklist

What I've actually built (not just claimed):

- [x] **JWT secret validation** — server refuses to start if `jwt_secret` is missing or under 32 chars.
- [x] **Env-var secrets override** — `JWT_SECRET`, `DB_PASSWORD` etc. take precedence over `config.ini`.
- [x] **Graceful shutdown** — SIGTERM triggers `srv.Shutdown(ctx)` with a 15s drain.
- [x] **Real readiness probe** — `/healthz` pings the DB; returns 503 on failure so a load balancer pulls the pod out.
- [x] **Rate limiting** — token-bucket per-IP on `/auth/*` (5/min) to slow credential stuffing.
- [x] **Security headers** — X-Content-Type-Options, X-Frame-Options, Referrer-Policy, X-XSS-Protection.
- [x] **Structured logs with request IDs** — `slog` JSON, request ID propagated via `X-Request-ID`.
- [x] **No GORM error leaks** — services classify errors, controllers map; only safe messages reach clients. Real errors are logged.
- [x] **Soft deletes** — GORM's `DeletedAt`; second-delete returns 404.
- [x] **Self-protection on admin actions** — admin can't demote / deactivate / delete their own account.
- [x] **CI on every push** — backend (build, vet, test with `-race`) + frontend (lint, type-check, test, build).
- [x] **Reversible migrations** with `Down` and a CLI rollback.
- [x] **Test coverage on critical paths** — auth, task CRUD, admin operations, error envelope.

What's intentionally **not** in this iteration (with the rationale):

- ❌ **Multi-tenancy / organizations** — see [ADR 005](docs/decisions/005-multi-tenancy-roadmap.md). It's the next big change, big enough to be its own project.
- ❌ **Refresh tokens** — current tokens expire in 3 days. Switching to short-lived access + httpOnly refresh is the right next step but needs the cookie path nailed down.
- ❌ **Email flows** (verification, password reset) — needs a transactional email provider (Resend / Postmark). Plumbing is in place; provider is the missing piece.
- ❌ **Stripe / billing** — same reason; needs a real account.
- ❌ **Distributed rate limiting** — the in-memory limiter is per-process. Multi-replica deploys need Redis (noted in the rate-limit middleware).

---

## Deploying for free

See [`docs/DEPLOY.md`](docs/DEPLOY.md) for a step-by-step guide using
**Vercel** (frontend) + **Render** (backend) + **Neon** (Postgres).
All free tier, no credit card required.

---

## Decisions worth reading

These ADRs explain the choices that aren't obvious from the code:

- [001 — Versioned raw-SQL migrations over GORM AutoMigrate](docs/decisions/001-raw-sql-migrations.md)
- [002 — AppError envelope and HTTP error responses](docs/decisions/002-apperror-envelope.md)
- [003 — TanStack Query key factories live with their service](docs/decisions/003-query-key-factories.md)
- [004 — Test layout: co-located backend, mirrored frontend](docs/decisions/004-test-layout.md)
- [005 — Multi-tenancy is the next big change (proposed)](docs/decisions/005-multi-tenancy-roadmap.md)
- [006 — Frontend organized by layer, not by feature](docs/decisions/006-frontend-layer-layout.md)

---

## License

MIT. Built as a portfolio project — feel free to use any of the patterns.
