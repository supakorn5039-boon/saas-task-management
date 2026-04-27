# 003 — Query-key factories live with their service

**Status:** Accepted &nbsp;·&nbsp; **Date:** 2026-04-24

## Context

TanStack Query's cache is keyed by a stable array. Every query and every
mutation that wants to invalidate something has to agree on what those keys
look like. The wrong shape:

```ts
useQuery({ queryKey: ["tasks", "list", page], ... })   // route 1
useQuery({ queryKey: ["task", "list", page], ... })    // route 2 (typo!)
```

Two queries, two caches, mutations invalidate the wrong one, list doesn't
refetch after a delete. This is the most common React Query bug in the wild.

## Decision

Every service exports a `*Keys` object whose shape is the source of truth for
that domain:

```ts
// services/task.service.ts
export const taskKeys = {
  all:    ["tasks"] as const,
  lists:  () => [...taskKeys.all, "list"] as const,
  list:   (params: ListTasksParams) => [...taskKeys.lists(), params] as const,
  details:() => [...taskKeys.all, "detail"] as const,
  detail: (id: number) => [...taskKeys.details(), id] as const,
};
```

Hooks **import** the factory; they don't write key arrays inline. Mutations
invalidate via `qc.invalidateQueries({ queryKey: taskKeys.all })` — one call
hits every list and detail because they all share the prefix.

Same pattern in `auth.service.ts` (`authKeys`) and `admin.service.ts`
(`adminKeys`).

## Consequences

**Good**

- Renaming or restructuring a key means changing it in one place.
- Invalidation is hierarchical: invalidate "tasks" and every list/detail under
  it gets refetched. Invalidate "tasks/list" and details are untouched.
- Routes and components don't need to know React Query internals — they just
  call the hooks.

**Trade-offs**

- Slightly more code than inline arrays. Worth it the first time you grep for
  "tasks" trying to find a stale cache key.
- The factory has to be invoked (`taskKeys.list(...)`) — passing
  `taskKeys.list` directly won't work. It's a small footgun.

## Notes for the next maintainer

If you're tempted to write `["tasks", "something"]` directly in a hook, add a
factory method instead. Keys are public API of the service — treat them that way.
