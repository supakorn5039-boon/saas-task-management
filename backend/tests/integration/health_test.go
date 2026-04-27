package integration

import (
	"net/http"
	"strings"
	"testing"
)

func TestHealth_Ping(t *testing.T) {
	srv := newServer(t)
	w := srv.Do(http.MethodGet, "/api/ping", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "pong") {
		t.Errorf("body = %s", w.Body.String())
	}
}

func TestHealth_Healthz(t *testing.T) {
	srv := newServer(t)
	w := srv.Do(http.MethodGet, "/api/healthz", nil, "")
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Errorf("body should be ok, got: %s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"db":"ok"`) {
		t.Errorf("body should report db ok, got: %s", w.Body.String())
	}
}
