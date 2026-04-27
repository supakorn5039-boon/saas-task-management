package integration

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestAuth_Register(t *testing.T) {
	srv := newServer(t)

	t.Run("success returns token + user", func(t *testing.T) {
		w := srv.Do(http.MethodPost, "/api/auth/register", map[string]string{
			"email":    "alice@example.com",
			"password": "password123",
		}, "")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
		var body struct {
			Token string                 `json:"token"`
			User  map[string]interface{} `json:"user"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if body.Token == "" {
			t.Error("expected non-empty token")
		}
		if body.User["email"] != "alice@example.com" {
			t.Errorf("user.email = %v", body.User["email"])
		}
	})

	t.Run("rejects short password (binding min=8)", func(t *testing.T) {
		w := srv.Do(http.MethodPost, "/api/auth/register", map[string]string{
			"email":    "bob@example.com",
			"password": "short",
		}, "")
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("rejects malformed email", func(t *testing.T) {
		w := srv.Do(http.MethodPost, "/api/auth/register", map[string]string{
			"email":    "not-an-email",
			"password": "password123",
		}, "")
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("duplicate email returns 409", func(t *testing.T) {
		w := srv.Do(http.MethodPost, "/api/auth/register", map[string]string{
			"email":    "alice@example.com",
			"password": "password123",
		}, "")
		if w.Code != http.StatusConflict {
			t.Errorf("status = %d, want 409", w.Code)
		}
	})
}

func TestAuth_Login(t *testing.T) {
	srv := newServer(t)
	srv.Register("login@example.com", "password123")

	t.Run("success", func(t *testing.T) {
		w := srv.Do(http.MethodPost, "/api/auth/login", map[string]string{
			"email":    "login@example.com",
			"password": "password123",
		}, "")
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
	})

	t.Run("wrong password returns 401 with generic message", func(t *testing.T) {
		w := srv.Do(http.MethodPost, "/api/auth/login", map[string]string{
			"email":    "login@example.com",
			"password": "wrong",
		}, "")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want 401", w.Code)
		}
		// The same generic message is used for both wrong-password and
		// unknown-email — so an attacker can't enumerate accounts.
		if !strings.Contains(w.Body.String(), "invalid email or password") {
			t.Errorf("body should hide which field is wrong, got: %s", w.Body.String())
		}
	})

	t.Run("unknown email returns identical 401", func(t *testing.T) {
		w := srv.Do(http.MethodPost, "/api/auth/login", map[string]string{
			"email":    "nobody@example.com",
			"password": "whatever123",
		}, "")
		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want 401", w.Code)
		}
		if !strings.Contains(w.Body.String(), "invalid email or password") {
			t.Errorf("body should match wrong-password text, got: %s", w.Body.String())
		}
	})
}

func TestAuth_RateLimit(t *testing.T) {
	srv := newServer(t)

	// The rate limiter is /auth/* — 5 requests per IP per minute. Register
	// itself counts toward that budget, so after registering we have 4
	// successful logins available before the limiter kicks in.
	srv.Register("rl@example.com", "password123") // budget used: 1/5

	for i := range 4 {
		w := srv.Do(http.MethodPost, "/api/auth/login", map[string]string{
			"email":    "rl@example.com",
			"password": "password123",
		}, "")
		if w.Code != http.StatusOK {
			t.Fatalf("login attempt %d (inside window): status = %d", i+1, w.Code)
		}
	}
	// Next call is the 6th /auth/* request → 429.
	w := srv.Do(http.MethodPost, "/api/auth/login", map[string]string{
		"email":    "rl@example.com",
		"password": "password123",
	}, "")
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("over window: status = %d, want 429", w.Code)
	}
	if w.Header().Get("Retry-After") == "" {
		t.Error("429 should include Retry-After header")
	}
}

// SecurityHeaders is a single sweep — verifying the headers are set on a
// representative response is enough; the middleware applies the same headers
// to every response.
func TestSecurityHeadersOnResponses(t *testing.T) {
	srv := newServer(t)
	w := srv.Do(http.MethodGet, "/api/ping", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}

	cases := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"X-XSS-Protection":       "0",
	}
	for header, want := range cases {
		if got := w.Header().Get(header); got != want {
			t.Errorf("%s = %q, want %q", header, got, want)
		}
	}
}

// RequestID middleware runs ahead of everything — the response should always
// carry an X-Request-ID, even on a no-handler endpoint.
func TestRequestIDIsAttached(t *testing.T) {
	srv := newServer(t)
	w := srv.Do(http.MethodGet, "/api/ping", nil, "")
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("response missing X-Request-ID header")
	}
}
