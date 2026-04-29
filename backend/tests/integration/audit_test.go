package integration

import (
	"net/http"
	"testing"
)

// TestAudit_LoginRecorded — login success creates a row in audit_logs that
// the actor can see in /user/activity, and the admin can see in
// /admin/audit-logs.
func TestAudit_LoginRecorded(t *testing.T) {
	srv := newServer(t)
	srv.Register("admin@example.com", "password123")
	srv.AsAdmin("admin@example.com")
	adminToken := srv.Login("admin@example.com", "password123")

	w := srv.Do(http.MethodGet, "/api/admin/audit-logs?action=auth.login", nil, adminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	var body struct {
		Data []struct {
			Action     string `json:"action"`
			ActorEmail string `json:"actorEmail"`
			Status     string `json:"status"`
		} `json:"data"`
		Meta struct{ Total int }
	}
	mustUnmarshal(t, w.Body.Bytes(), &body)
	if body.Meta.Total < 1 {
		t.Fatalf("expected at least 1 login event, got %d", body.Meta.Total)
	}
	for _, e := range body.Data {
		if e.Action != "auth.login" {
			t.Errorf("action = %q, want auth.login", e.Action)
		}
		if e.Status != "success" {
			t.Errorf("status = %q, want success", e.Status)
		}
	}
}

func TestAudit_FailedLoginRecorded(t *testing.T) {
	srv := newServer(t)
	srv.Register("user@example.com", "password123")

	// Bad password — should produce an auth.login_failed entry.
	w := srv.Do(http.MethodPost, "/api/auth/login",
		map[string]string{"email": "user@example.com", "password": "wrong-password"}, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}

	srv.AsAdmin("user@example.com")
	adminToken := srv.Login("user@example.com", "password123")

	w = srv.Do(http.MethodGet, "/api/admin/audit-logs?action=auth.login_failed", nil, adminToken)
	var body struct {
		Meta struct{ Total int }
	}
	mustUnmarshal(t, w.Body.Bytes(), &body)
	if body.Meta.Total != 1 {
		t.Errorf("Total = %d, want 1 failed-login event", body.Meta.Total)
	}
}

func TestAudit_MyActivityScopedToActor(t *testing.T) {
	srv := newServer(t)
	aliceToken, _, _ := srv.Register("alice@example.com", "password123")
	srv.Register("bob@example.com", "password123")
	srv.Login("bob@example.com", "password123") // creates an audit row for bob

	// Alice's /user/activity should only contain alice's events.
	w := srv.Do(http.MethodGet, "/api/user/activity", nil, aliceToken)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	var body struct {
		Data []struct {
			ActorEmail string `json:"actorEmail"`
		} `json:"data"`
	}
	mustUnmarshal(t, w.Body.Bytes(), &body)
	for _, e := range body.Data {
		if e.ActorEmail != "alice@example.com" {
			t.Errorf("found bob's event in alice's activity feed: %+v", e)
		}
	}
}

func TestAudit_AdminEndpointBlockedForRegularUser(t *testing.T) {
	srv := newServer(t)
	userToken, _, _ := srv.Register("user@example.com", "password123")

	w := srv.Do(http.MethodGet, "/api/admin/audit-logs", nil, userToken)
	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", w.Code)
	}
}

func TestAudit_LogoutRecorded(t *testing.T) {
	srv := newServer(t)
	token, _, _ := srv.Register("logout@example.com", "password123")

	w := srv.Do(http.MethodPost, "/api/auth/logout", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	// Promote and check the audit log records the logout.
	srv.AsAdmin("logout@example.com")
	adminToken := srv.Login("logout@example.com", "password123")
	w = srv.Do(http.MethodGet, "/api/admin/audit-logs?action=auth.logout", nil, adminToken)
	var body struct {
		Meta struct{ Total int }
	}
	mustUnmarshal(t, w.Body.Bytes(), &body)
	if body.Meta.Total != 1 {
		t.Errorf("Total = %d, want 1 logout event", body.Meta.Total)
	}
}
