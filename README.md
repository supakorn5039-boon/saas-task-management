# SaaS Task Management

A full-stack task management application built with **Go + React**, featuring JWT authentication, role-based access, an admin dashboard, and a colorful modern UI.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?logo=typescript&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.5-4169E1?logo=postgresql&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-4.2-06B6D4?logo=tailwindcss&logoColor=white)

---

## Tech Stack

| Layer        | Technology                                                                |
| ------------ | ------------------------------------------------------------------------- |
| **Frontend** | React 19, TypeScript, TanStack Router / Query / Table, Zustand            |
| **Backend**  | Go 1.25, Gin, GORM, JWT                                                   |
| **Database** | PostgreSQL 17.5                                                           |
| **Styling**  | Tailwind CSS 4, shadcn/ui, Lucide Icons                                   |
| **Tooling**  | Vite 8, Lefthook, ESLint, Prettier, Docker                                |

---

## Features

### App
- **Authentication** — Register, login, JWT-protected routes, persistent auth via Zustand
- **Tasks** — Create with title + description, three-state status (`todo` / `in_progress` / `done`), delete
- **Paginated data table** — Server-side sort, filter by status, full-text search, page-size selector (10/20/50/100)
- **Status filter tabs** with live counts per status (one round-trip)
- **Admin sidebar shell** — Collapsible nav with role-aware menu items (admin-only sections hidden for regular users)
- **Dashboard** — KPI cards (total / todo / in progress / done) + recent activity, all from a single API call
- **Theme** — Vibrant indigo brand, status-colored badges (slate / amber / emerald), tooltips on hover for action buttons
- **Toast notifications** — Sonner for success/error feedback

### Architecture
- **Migration framework** — Versioned, immutable raw-SQL migrations with auto-run on boot + rollback support
- **Seeder framework** — Registry pattern via `init()`, GORM with `clause.OnConflict` for idempotency
- **Query key factory** — TanStack Query keys live with the service; hooks consume them so routes don't know React Query internals
- **Clean separation** — `types/` (pure), `constants/` (display labels), `styles/` (Tailwind), `components/` (UI), `hooks/` (data), `services/` (API)
- **Reusable data-table primitives** — `DataTablePagination` and `DataTableSortHeader` work for any future table

---

## Project Structure

```
saas-task-management/
├── backend/
│   ├── src/
│   │   ├── apiwebserver/
│   │   │   ├── controller/      # HTTP handlers
│   │   │   ├── service/         # Business logic + GORM queries
│   │   │   └── middleware/      # JWT auth
│   │   ├── config/              # INI config loader
│   │   ├── database/
│   │   │   ├── database.go      # Connect + auto-run migrations
│   │   │   ├── model/           # GORM models + DTOs
│   │   │   ├── migration/       # Versioned raw-SQL migrations
│   │   │   └── seeder/          # Registry-based seeders
│   │   ├── pkg/                 # API server bootstrap (CORS, routes)
│   │   ├── security/            # JWT + bcrypt helpers
│   │   └── main.go              # Entry point + CLI commands
│   ├── bruno/                   # API collection
│   ├── docker-compose.yml
│   └── config.ini               # (gitignored — copy from config.example.ini)
│
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   │   ├── ui/              # shadcn/ui primitives
│   │   │   ├── layout/          # AppShell, AppSidebar, UserMenu, nav-config
│   │   │   ├── data-table/      # Reusable pagination + sort header
│   │   │   ├── tasks/           # TaskTable, TaskStatusBadge, TaskStatusSelect
│   │   │   └── dashboard/       # KpiCard
│   │   ├── routes/              # File-based pages (TanStack Router)
│   │   ├── services/            # API clients + query-key factories
│   │   ├── hooks/               # useTasks + mutation hooks
│   │   ├── store/               # Zustand auth store
│   │   ├── types/               # Pure TypeScript types
│   │   ├── constants/           # Display labels (TASK_STATUS_LABEL, etc.)
│   │   ├── styles/              # Tailwind class maps per domain
│   │   └── lib/                 # Utilities + axios instance
│   ├── package.json
│   └── vite.config.ts
│
├── lefthook.yml                 # Pre-commit / pre-push hook config
└── README.md
```

---

## API Endpoints

### Public

| Method | Endpoint             | Description       |
| ------ | -------------------- | ----------------- |
| `GET`  | `/api/ping`          | Health check      |
| `POST` | `/api/auth/login`    | User login        |
| `POST` | `/api/auth/register` | User registration |

### Protected (JWT required)

| Method   | Endpoint            | Description                                              |
| -------- | ------------------- | -------------------------------------------------------- |
| `GET`    | `/api/user/profile` | Get current user                                         |
| `GET`    | `/api/tasks`        | List tasks (paginated + filterable — see query params)   |
| `POST`   | `/api/tasks`        | Create a task                                            |
| `PUT`    | `/api/tasks/:id`    | Update task status                                       |
| `DELETE` | `/api/tasks/:id`    | Delete a task                                            |

### `GET /api/tasks` query params

| Param      | Default      | Notes                                                            |
| ---------- | ------------ | ---------------------------------------------------------------- |
| `page`     | `1`          | Page number, ≥ 1                                                 |
| `per_page` | `10`         | Page size, clamped 1–100                                         |
| `status`   | (none)       | Filter: `todo` \| `in_progress` \| `done`                        |
| `search`   | (none)       | Substring match on `title` and `description` (case-insensitive)  |
| `sort`     | `created_at` | One of `created_at`, `updated_at`, `title`, `status` (whitelist) |
| `order`    | `desc`       | `asc` or `desc`                                                  |

Response shape:

```json
{
  "data":  [ { "id": 1, "title": "...", "description": "...", "status": "todo", "createdAt": "..." } ],
  "meta":  { "page": 1, "perPage": 10, "total": 47 },
  "counts": { "all": 47, "todo": 30, "in_progress": 10, "done": 7 }
}
```

---

## Getting Started

### Prerequisites

- **Go** 1.25+
- **Node.js** 18+
- **Docker** (for PostgreSQL)
- **Lefthook** (`brew install lefthook` on macOS)

### 1. Start the Database

```bash
cd backend
docker-compose up -d
```

### 2. Configure the Backend

```bash
cp backend/config.example.ini backend/config.ini
# edit jwt_secret + DB credentials as needed
```

### 3. Run the Backend

```bash
cd backend
go run ./src
```

> Backend runs on `http://localhost:8080`.
> Migrations auto-run on every boot via `database.Connect()`.

CLI commands:

```bash
go run ./src                  # start the HTTP server
go run ./src seed             # seed default roles + users (admin@example.com / user@example.com, password: password123)
go run ./src migrate:status   # show migration status
go run ./src migrate:rollback # rollback the last migration
```

### 4. Run the Frontend

```bash
cd frontend
npm install
npm run dev
```

> Frontend runs on `http://localhost:5173`.

### 5. Enable git hooks (once per clone)

```bash
lefthook install
```

Wires up `.git/hooks/pre-commit` and `.git/hooks/pre-push` so every commit/push runs the right checks (see [Git Hooks](#git-hooks) below).

### 6. (Optional) Use a custom local domain

Instead of `localhost:5173`, run the app at `http://saas-management.local`:

```bash
echo "127.0.0.1 saas-management.local" | sudo tee -a /etc/hosts
sudo brew services start nginx   # uses /opt/homebrew/etc/nginx/servers/saas-management.conf
```

Nginx routes `/api/*` to the Go backend (`:8080`) and everything else to the Vite dev server (`:5173`), with WebSocket upgrade for HMR. The frontend is whitelisted via `vite.config.ts` `server.allowedHosts`.

---

## Git Hooks

Managed by [lefthook](https://github.com/evilmartians/lefthook) — single Go binary, configured in `lefthook.yml`.

### Pre-commit (only on changed files, parallel ~3s)

| Match                          | Check                                |
| ------------------------------ | ------------------------------------ |
| `frontend/src/**/*.{ts,tsx}`   | eslint --fix, prettier, tsc          |
| `backend/**/*.go`              | go build, go vet, go test            |

### Pre-push (heavier checks)

| Check        | What it does                                  |
| ------------ | --------------------------------------------- |
| vite-build   | Full production build of the frontend         |
| todo-scan    | Reports `TODO` / `FIXME` markers              |
| secret-scan  | Blocks pushes with hardcoded secret patterns  |

Run any stage manually:

```bash
lefthook run pre-commit --all-files
lefthook run pre-push --all-files
```

---

## Frontend Scripts

| Script             | Command                  |
| ------------------ | ------------------------ |
| `npm run dev`      | Start Vite dev server    |
| `npm run build`    | Production build         |
| `npm run preview`  | Preview production build |
| `npm run lint`     | ESLint                   |
| `npm run lint:fix` | Auto-fix lint issues     |
| `npm run format`   | Prettier write           |
| `npm run type-check` | TypeScript check       |

---

## Environment Variables

### Frontend

Create `frontend/.env` (optional — defaults work for local dev):

```
VITE_API_URL="http://localhost:8080/api"
```

### Backend

Configure `backend/config.ini`:

```ini
[server]
port       = 8080
production = false
jwt_secret = your-secret-key

[database]
host     = localhost
user     = username
password = password
database = postgres
port     = 5432
```

---

## Default Seeded Users

Run `go run ./src seed` to create:

| Email                | Role  | Password      |
| -------------------- | ----- | ------------- |
| `admin@example.com`  | admin | `password123` |
| `user@example.com`   | user  | `password123` |

Roles control sidebar visibility — `admin` sees the **Users** and **Settings** menu items; `user` does not.

---

## License

This project is for educational and portfolio purposes.
