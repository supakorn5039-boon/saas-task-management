# 005 — Multi-tenancy is the next big change

**Status:** Proposed &nbsp;·&nbsp; **Date:** 2026-04-27

## Context

The product today is a single-tenant task tracker: every task belongs to a
user, and every user is in the same global namespace. This is fine for a
portfolio MVP and for a personal tool — but it's not SaaS in the meaningful
sense of the word, and the entire data model would shift the moment we tried
to add an "organization" or "team" concept.

This ADR captures **the plan**, before the work starts, so that later code
review has a reference and so the reasoning is visible to anyone reading the
repo.

## Proposed shape

### Domain

- New `organizations` table — each org has a name, slug, owner.
- New `organization_members` join table — `(organization_id, user_id, role,
  invited_at, joined_at)`. The existing `users.role` becomes a *workspace*
  role (admin/manager/user) and stays. The new `organization_members.role`
  becomes the *org* role (owner/admin/member).
- Every tenant-scoped table gains an `organization_id` foreign key:
  `tasks` first, then anything we add later.

### Auth flow

- Login still returns a JWT; the JWT now also carries `current_organization_id`.
- Users with multiple orgs get an org-switcher in the UI; the switcher updates
  the JWT (re-issues a token).
- Middleware reads `current_organization_id` from the claims and sets it on
  the request context. **Every** service query must filter by it — this is
  the line we don't cross.

### Admin scope

- Today's "admin" is global. After this change, admins are scoped to an org.
- A super-admin (system role) is added for support cases — visible in the
  same `users` table by `system_role = 'super_admin'`. Most admin code never
  sees this.

### API impact

- Routes stay the same; `organization_id` is implicit in the JWT, not in the
  URL. Tradeoff: URLs are stable, but tenant context is invisible at the
  request line.
- Tasks, users, settings — every endpoint that currently filters by
  `user_id` also filters by `organization_id`.

## Why a separate ADR before any code lands

Multi-tenancy touches every model, every query, every test, every controller.
Doing this badly creates a class of bug — *cross-tenant data leak* — that
silently violates customer expectations and is hard to audit. Pinning the
plan first lets us:

1. Review the model on paper (ADR feedback) before committing to migrations.
2. Set the rule once: **every tenant-scoped query MUST filter by
   `organization_id`**, enforced via service-layer helpers, never trusted to
   the controller.
3. Start tests for tenancy isolation on day one. Each cross-tenant assertion
   becomes a regression test.

## Consequences (when implemented)

**Good**

- The product becomes a real SaaS — invitable, billable per-org, isolatable.
- Every existing feature gets multi-user collaboration "for free" once the
  org model lands (assignees, shared tasks, shared dashboards).

**Trade-offs**

- Every existing service test needs an org fixture. Boilerplate cost.
- The "single user, my private tasks" experience becomes "your default
  organization" — UX work to make the first-time login feel as smooth as it
  does now.
- The migration that adds `organization_id` to existing rows has to backfill
  cleanly — design before writing.

## Status

This ADR is **Proposed**, not Accepted, until the migration plan above has
been written out in detail (a follow-up doc) and reviewed.
