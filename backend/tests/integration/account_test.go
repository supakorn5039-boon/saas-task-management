package integration

import (
	"net/http"
	"strings"
	"testing"
)

func TestAccount_ChangePassword(t *testing.T) {
	srv := newServer(t)
	// One register call (budget: 1/5). Reuse the JWT — it doesn't get
	// invalidated by password changes (signed with a static secret), so we
	// can do every subtest without burning more /auth/* budget.
	token, _, _ := srv.Register("pw@example.com", "password123")

	t.Run("rejects same-as-current new password", func(t *testing.T) {
		w := srv.Do(http.MethodPut, "/api/user/password", map[string]string{
			"currentPassword": "password123",
			"newPassword":     "password123",
		}, token)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
		if !strings.Contains(w.Body.String(), "different") {
			t.Errorf("expected 'different' hint, got: %s", w.Body.String())
		}
	})

	t.Run("rejects wrong current password", func(t *testing.T) {
		w := srv.Do(http.MethodPut, "/api/user/password", map[string]string{
			"currentPassword": "wrong",
			"newPassword":     "another-strong-pass",
		}, token)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want 401", w.Code)
		}
	})

	t.Run("rejects too-short new password", func(t *testing.T) {
		w := srv.Do(http.MethodPut, "/api/user/password", map[string]string{
			"currentPassword": "password123",
			"newPassword":     "short",
		}, token)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("happy path", func(t *testing.T) {
		w := srv.Do(http.MethodPut, "/api/user/password", map[string]string{
			"currentPassword": "password123",
			"newPassword":     "another-strong-pass",
		}, token)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}

		// Old password no longer works (login attempt #2 — within budget).
		w = srv.Do(http.MethodPost, "/api/auth/login", map[string]string{
			"email":    "pw@example.com",
			"password": "password123",
		}, "")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("login with old password: status = %d, want 401", w.Code)
		}

		// New password works (login attempt #3).
		w = srv.Do(http.MethodPost, "/api/auth/login", map[string]string{
			"email":    "pw@example.com",
			"password": "another-strong-pass",
		}, "")
		if w.Code != http.StatusOK {
			t.Errorf("login with new password: status = %d, want 200", w.Code)
		}
	})
}

func TestAccount_GetProfile(t *testing.T) {
	srv := newServer(t)
	token, _, _ := srv.Register("profile@example.com", "password123")

	w := srv.Do(http.MethodGet, "/api/user/profile", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "profile@example.com") {
		t.Errorf("body should include email, got: %s", w.Body.String())
	}
}
