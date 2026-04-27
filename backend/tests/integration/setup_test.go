// Package integration holds black-box HTTP tests that exercise the real
// Gin engine end-to-end. They go through the same middleware (auth, RBAC,
// rate limit, security headers, request logger) and the same routes that
// production uses, so they catch wiring bugs that pure service-level tests
// can't.
//
// Why this directory exists separately from the unit tests:
//   - Unit tests (in src/apiwebserver/service/, src/apperror/) must live in
//     the same package as the code under test — Go's privacy model gives
//     package-level visibility, so tests in another directory can't see
//     internal symbols. That's a language rule, not a convention.
//   - These integration tests touch only the public HTTP surface, so they're
//     free to live wherever. Pulling them out keeps the source tree clean
//     and matches the layout used by larger Go projects (kubernetes/, etcd/,
//     etc).
package integration

import (
	"github.com/supakorn5039-boon/saas-task-backend/src/testhelpers"
	"testing"
)

// newServer is a thin convenience wrapper so individual tests don't have to
// know about testhelpers — they can just call `srv := newServer(t)`.
func newServer(t *testing.T) *testhelpers.TestServer {
	t.Helper()
	return testhelpers.NewTestServer(t)
}
