# Architecture Decision Records

Short, dated records of the non-obvious choices in this codebase. Each one
describes the **context**, the **decision**, and the **consequences** — not
just what was done, but why. New ADRs go on the end with the next number.

| #   | Title                                                | Status     |
| --- | ---------------------------------------------------- | ---------- |
| 001 | [Versioned raw-SQL migrations over GORM AutoMigrate](001-raw-sql-migrations.md) | Accepted   |
| 002 | [AppError + sentinel handling for HTTP responses](002-apperror-envelope.md)     | Accepted   |
| 003 | [TanStack Query with co-located key factories](003-query-key-factories.md)     | Accepted   |
| 004 | [Co-locate tests next to code (backend) / `__tests__/` mirror (frontend)](004-test-layout.md) | Accepted |
| 005 | [Multi-tenancy is the next big change](005-multi-tenancy-roadmap.md)            | Proposed   |
| 006 | [Frontend organized by layer, not by feature](006-frontend-layer-layout.md)     | Accepted   |
