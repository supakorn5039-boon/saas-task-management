# 006 — Frontend organized by layer, not by feature

**Status:** Accepted &nbsp;·&nbsp; **Date:** 2026-04-27

## Context

Two common ways to organize a React app's `src/`:

**Layer-based** — folders are *kinds* of file:

```
src/
  services/        # all API clients
  hooks/           # all hooks
  validators/      # all Zod schemas
  types/           # all shared types
  components/
    ui/            # shadcn primitives
    shared/        # cross-feature reusable
    tasks/        # task-specific
    admin/        # admin-specific
    layout/       # app shell
  routes/         # TanStack Router file-based pages
  store/          # Zustand stores
  styles/         # Tailwind class maps
  constants/      # display labels, etc.
  lib/            # cross-cutting utilities
```

**Feature-sliced** — folders are *features*:

```
src/
  features/
    auth/{api.ts, validators.ts, types.ts}
    tasks/{api.ts, validators.ts, types.ts, hooks/, components/}
    admin/{...}
  components/{ui, shared, layout}
  routes/
  lib/
```

Both work. They're mutually exclusive — you can't really have both without
the project feeling inconsistent.

## Decision

We use the **layer-based** layout.

A short experimental move toward feature-sliced was reverted before any
inconsistency could be merged, so the codebase is fully one or the other.

## Why layer-based for this project

- **Project size.** The product has 5 domains (auth, tasks, admin, account,
  dashboard). Layer-based reads cleanly at this size; feature-sliced doesn't
  start outscoring it until 8–10 features.
- **Consistency win.** "All my Zod schemas live in `validators/`" is a single
  rule that's easy for any contributor (and any reviewer) to keep in mind.
- **Low cognitive overhead.** A new contributor can find any file in two
  hops: pick the kind (`services/`), pick the feature (`task.service.ts`).
- **Routes-as-pages.** TanStack Router uses file-based routing, so the URL
  tree already lives in `routes/`. We don't need a parallel `features/`
  hierarchy to express "this is a page."

## Consequences

**Good**

- Predictable layout. Once you know the conventions for one layer, the rest
  follow the same shape.
- Refactors stay in-place — moving a file rarely needs an import path update
  beyond what the IDE handles.

**Trade-offs**

- "Show me everything related to admin" requires touching 5 folders. We
  accept this — it's the explicit cost of grouping by kind.
- If the project grows past ~10 features, the top-level folders get crowded.
  At that point, switching to feature-sliced should be revisited (and a new
  ADR written explaining the move).

## Rule of thumb for placement

| File                    | Where it goes                       |
| ----------------------- | ----------------------------------- |
| API client              | `services/<feature>.service.ts`     |
| TanStack Query hook     | `hooks/use-<feature>.ts`            |
| Zod schema              | `validators/<feature>.validator.ts` |
| Cross-feature TS type   | `types/<feature>.ts`                |
| Feature-only component  | `components/<feature>/<name>.tsx`   |
| shadcn primitive        | `components/ui/<name>.tsx`          |
| Cross-feature reusable  | `components/shared/<name>.tsx`      |
| Page                    | `routes/<path>.lazy.tsx`            |

## Notes for the next maintainer

If you find yourself adding a top-level folder (e.g. `repositories/`, `dto/`),
ask whether it really earns its keep. The current set already covers most
shapes a React file can take. New layers are sometimes warranted — but each
new top-level folder is a small tax on every contributor's mental model of
the repo, so add them with intent.
