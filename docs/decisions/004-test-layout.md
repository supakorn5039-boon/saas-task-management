# 004 — Test layout: backend co-located, frontend in `__tests__/`

**Status:** Accepted &nbsp;·&nbsp; **Date:** 2026-04-27

## Context

Two languages, two ecosystems, two conventions. This ADR exists so the next
person doesn't ask "why is the layout different?" — there's a real reason.

## Decision

### Backend (Go) — co-located, mandatory

Go's testing tooling requires test files to live in the same directory (and
typically the same package or a `_test`-suffixed sibling package) as the code
under test:

```
backend/src/apiwebserver/service/
  task_service.go
  task_service_test.go
  auth_service.go
  auth_service_test.go
```

There is no opt-out — `go test ./...` discovers tests this way.

### Frontend (Vitest) — separate `__tests__/` tree

Frontend tests mirror the source tree under `frontend/src/__tests__/`:

```
src/
  hooks/
    use-task-list-state.ts
  __tests__/
    hooks/
      use-task-list-state.test.ts
```

`vitest.config.ts` has `include: ["src/__tests__/**/*.{test,spec}.{ts,tsx}"]`.
Tests import via the `@/` alias (`@/hooks/use-task-list-state`) so they don't
care where they live in the tree.

## Why different on the frontend?

Modern Vitest/Jest projects more commonly co-locate (`foo.ts` next to
`foo.test.ts`). We considered both:

| Co-located                                  | Separate tree (chosen)                  |
| ------------------------------------------- | --------------------------------------- |
| Source folder shows tests too               | Source folder is "shippable code only"  |
| Easy to grep "for the test of X"            | One place to look at the test surface   |
| Renames are atomic                          | Renames need to update both trees       |
| Vitest default                              | Slightly more config                    |

The separate tree was chosen because the source folder reads cleaner when
you're scanning the actual product code — tests don't dilute the file list of
`hooks/` or `lib/`. Imports use the `@/` alias, so refactor cost is minimal.

## Consequences

**Good**

- Production source folders contain only production code.
- Test layout is self-documenting — the mirror tells you what's covered.

**Trade-offs**

- A renamed source file leaves an orphan test until you update the mirror.
  Mitigated by ESLint catching the broken import on next lint.
- Some IDE plugins assume co-location for "go to test" navigation.

## Notes for the next maintainer

If you ever want to switch to co-location, change `vitest.config.ts`'s
`include` glob to `src/**/*.{test,spec}.{ts,tsx}` and move the files. The
imports already use `@/` aliases, so nothing else has to change.
