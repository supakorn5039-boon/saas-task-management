# SaaS Task Management

A full-stack task management application built with **Go + React**, featuring JWT authentication, role-based access, and a modern UI.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)
![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?logo=typescript&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.5-4169E1?logo=postgresql&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-4.2-06B6D4?logo=tailwindcss&logoColor=white)

---

## Tech Stack

| Layer        | Technology                                                    |
| ------------ | ------------------------------------------------------------- |
| **Frontend** | React 19, TypeScript, TanStack Router, TanStack Query, Zustand |
| **Backend**  | Go, Gin, GORM, JWT                                           |
| **Database** | PostgreSQL 17.5                                               |
| **Styling**  | Tailwind CSS 4, shadcn/ui, Lucide Icons                      |
| **Tooling**  | Vite 8, Docker, Biome                                        |

---

## Features

- **Authentication** — Register, login, and JWT-protected routes
- **Task Management** — Create, complete, and delete tasks
- **User Profile** — View account info with auto-generated avatar
- **Dark Mode** — Theme switching with `next-themes`
- **Form Validation** — Zod schemas + React Hook Form
- **Toast Notifications** — Feedback via Sonner
- **Persistent Auth** — Zustand store with localStorage

---

## Project Structure

```
saas-task-management/
├── backend/
│   ├── src/
│   │   ├── app/              # Application bootstrap
│   │   ├── cmd/              # CLI commands (migration)
│   │   ├── config/           # Config loader (INI)
│   │   ├── database/         # DB connection & migrations
│   │   ├── models/           # Data models & DTOs
│   │   ├── security/         # JWT & bcrypt helpers
│   │   ├── server/
│   │   │   ├── controllers/  # Request handlers
│   │   │   ├── services/     # Business logic
│   │   │   ├── routes/       # Route definitions
│   │   │   └── middleware/   # Auth middleware
│   │   └── utils/            # Utility functions
│   ├── bruno/                # API collection (Bruno)
│   ├── docker/               # Docker configs
│   ├── docker-compose.yml
│   ├── Dockerfile
│   └── main.go               # Entry point
│
├── frontend/
│   ├── src/
│   │   ├── components/ui/    # shadcn/ui components
│   │   ├── routes/           # File-based pages
│   │   ├── services/         # API clients (Axios)
│   │   ├── store/            # Zustand auth store
│   │   ├── types/            # TypeScript types
│   │   └── lib/              # Utilities
│   ├── package.json
│   └── vite.config.ts
```

---

## API Endpoints

### Public

| Method | Endpoint          | Description       |
| ------ | ----------------- | ----------------- |
| `GET`  | `/api/ping`       | Health check      |
| `POST` | `/api/auth/login` | User login        |
| `POST` | `/api/auth/register` | User registration |

### Protected (JWT Required)

| Method   | Endpoint           | Description          |
| -------- | ------------------ | -------------------- |
| `GET`    | `/api/user/profile`| Get current user     |
| `GET`    | `/api/tasks`       | List all tasks       |
| `POST`   | `/api/tasks`       | Create a task        |
| `PATCH`  | `/api/tasks/:id`   | Toggle task status   |
| `DELETE` | `/api/tasks/:id`   | Delete a task        |

---

## Getting Started

### Prerequisites

- **Go** 1.25+
- **Node.js** 18+
- **Docker** (for PostgreSQL)

### 1. Start the Database

```bash
cd backend
docker-compose up -d
```

### 2. Run the Backend

```bash
cd backend
go run main.go
```

> Backend runs on `http://localhost:8080`

### 3. Run the Frontend

```bash
cd frontend
npm install
npm run dev
```

> Frontend runs on `http://localhost:5173`

---

## Scripts

### Frontend

| Script           | Command                        |
| ---------------- | ------------------------------ |
| `npm run dev`    | Start dev server               |
| `npm run build`  | Production build               |
| `npm run lint`   | Lint with Biome                |
| `npm run lint:fix`| Auto-fix lint issues          |
| `npm run format` | Format code with Biome         |
| `npm run preview`| Preview production build       |

---

## Environment Variables

### Frontend

Create `frontend/.env`:

```
VITE_API_URL="http://localhost:8080/api"
```

### Backend

Configure `backend/config.ini`:

```ini
[server]
port = 8080
production = false
jwt_secret = your-secret-key

[database]
host = localhost
user = username
password = password
database = postgres
port = 5432
```

---

## License

This project is for educational and portfolio purposes.
