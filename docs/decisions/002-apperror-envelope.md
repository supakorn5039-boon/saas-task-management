# 002 — AppError envelope and HTTP error responses

**Status:** Accepted &nbsp;·&nbsp; **Date:** 2026-04-26

## Context

Early versions of the API did this:

```go
if err != nil {
    errorResponse(c, err.Error(), http.StatusInternalServerError)
    return
}
```

Two problems:

1. **GORM strings leaked to clients.** A "duplicate key" Postgres error went
   straight to the browser. The client could learn the column name, the index
   name, sometimes the value being checked.
2. **Status codes were guessed at the controller.** A "record not found" from a
   service layer became a 500 because the controller didn't know to translate.

We needed every error path to either be safe-by-default or carry an explicit
status + client-safe message.

## Decision

A small `apperror` package defines:

```go
type AppError struct {
    Status  int
    Message string  // safe to send to clients
    Err     error   // wrapped cause, logged but never sent
}
```

with constructors `New`, `Wrap`, `NotFound`, `BadRequest`, `Unauthorized`,
`Conflict`, plus `FromGorm(err, msg)` that maps `gorm.ErrRecordNotFound` to a
404 and passes anything else through.

Services return `AppError` for any condition the client should be told about
(404, 401, 409, 400). Anything else is wrapped via `apperror.Wrap(err, 500, ...)`.

The shared controller helper:

```go
func errorResponse(c *gin.Context, err error) {
    if ae, ok := apperror.As(err); ok {
        // Log the underlying cause, send only the safe message.
        c.AbortWithStatusJSON(ae.Status, gin.H{"error": ae.Message})
        return
    }
    log.Printf("unhandled error: %v", err)
    c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
}
```

## Consequences

**Good**

- The client never sees an error string we didn't write ourselves.
- Status codes flow from the layer that knows them (the service), not the
  controller's guess.
- Untyped errors are loud (logged) and safe (always 500 with a generic message).
- The frontend can rely on the response always being `{ "error": "..." }` —
  see `frontend/src/lib/api-error.ts`.

**Trade-offs**

- Two layers of error handling: services return `*AppError`, controllers
  unwrap. Worth it; the alternative is leaks.
- Validation errors from Gin's binding (`binding:"required,min=8"`) bypass
  this — they go through `badRequest(c, err.Error())` directly because the
  validator's messages are already client-safe.

## Notes for the next maintainer

If you find yourself writing `err.Error()` into a `c.JSON` call, stop. Wrap
it as `apperror.Wrap(err, 500, "what happened")` and let the helper handle it.
