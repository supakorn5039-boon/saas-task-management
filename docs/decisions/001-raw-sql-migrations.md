# 001 — Versioned raw-SQL migrations over GORM AutoMigrate

**Status:** Accepted &nbsp;·&nbsp; **Date:** 2026-04-10

## Context

GORM ships with `AutoMigrate`, which inspects struct tags and `ALTER`s the
schema to match. It's seductive in dev — you change a struct, you restart, the
schema "just works" — but for any database that holds data you care about,
it's the wrong tool:

- **No history.** You can't tell what migrations have run on which environment.
- **No rollback.** Once a column is dropped, the original is gone.
- **No down path.** Rolling back a release means rolling forward into a fix.
- **Implicit ordering.** Two engineers add columns in different branches; the
  merge order in production is non-deterministic.

## Decision

Each schema change is a versioned, timestamped Go file that implements:

```go
type Migration interface {
    ID() string                     // matches the filename prefix
    Up(db *gorm.DB) error           // raw SQL
    Down(db *gorm.DB) error         // raw SQL — actual rollback, not a guess
}
```

Migrations are registered via `init()` and run on boot in ID order. A
`schema_migrations` table tracks which IDs have applied. `migrate:rollback` and
`migrate:status` CLI commands are exposed in `main.go`.

Files live in `backend/src/database/migration/` named
`YYYY_MM_DD_HH_MM_SS_<slug>.go`.

## Consequences

**Good**

- Rollbacks are real, not theoretical.
- Same SQL runs in dev and prod — no drift between AutoMigrate's interpretation
  of struct tags and the actual production DDL.
- A new env can be brought up with the same script as prod (run all migrations
  in order).

**Trade-offs**

- A schema change costs writing both `Up` and `Down`. We accept this — it's
  the minimum cost of a real rollback story.
- Migration files are immutable once shipped. Changing one rewrites history;
  fixes go forward as a new migration.

## Notes for the next maintainer

Don't add `db.AutoMigrate(...)` anywhere except test helpers (see
`backend/src/testhelpers/db.go`). Tests use AutoMigrate because they want a
fresh schema per run; production is the reverse.
