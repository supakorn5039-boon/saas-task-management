package testhelpers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/controller"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
)

// testJWTSecret is at least 32 chars to satisfy config-time validation.
const testJWTSecret = "this-is-a-test-jwt-secret-at-least-32-chars-long"

// TestServer wraps a Gin engine pre-wired with the same middleware and routes
// the production server uses. Each test gets its own engine so rate limiters
// and other in-memory state don't leak across tests.
type TestServer struct {
	t      *testing.T
	Engine *gin.Engine
}

// NewTestServer returns a fully-wired test server. SetupTestDB has already run
// against the test database, so services constructed via NewXService() use it.
func NewTestServer(t *testing.T) *TestServer {
	t.Helper()
	SetupTestDB(t)

	gin.SetMode(gin.TestMode)
	security.InitJWT(testJWTSecret)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())
	r.Use(middleware.SecurityHeaders())

	api := r.Group("/api")
	controller.NewAuthController().RegisterRoutes(api)
	controller.NewUserController().RegisterRoutes(api)
	controller.NewTaskController().RegisterRoutes(api)
	controller.NewAdminController().RegisterRoutes(api)

	// Mirror the inline /ping and /healthz routes from pkg.MountAPIWebServer
	// so integration tests can assert on them too. Keep behavior identical.
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	api.GET("/healthz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		sqlDB, err := database.DB.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "db": "unreachable"})
			return
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "db": "ping_failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "ok"})
	})

	return &TestServer{t: t, Engine: r}
}

// Do issues an in-process HTTP request against the test server. body may be
// any JSON-marshalable value (including nil for none); token is the bearer
// token (empty for unauth requests).
func (s *TestServer) Do(method, path string, body any, token string) *httptest.ResponseRecorder {
	s.t.Helper()

	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			s.t.Fatalf("marshal body: %v", err)
		}
		reader = bytes.NewReader(raw)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	s.Engine.ServeHTTP(w, req)
	return w
}

// Login registers a user and returns the issued JWT plus the user payload.
// The returned response shape mirrors the auth endpoint:
//
//	{ "token": "...", "user": { "id": ..., "email": "...", "role": "..." } }
//
// Convenient for tests that need an authenticated client.
func (s *TestServer) Register(email, password string) (token string, userID uint, role string) {
	s.t.Helper()
	w := s.Do(http.MethodPost, "/api/auth/register", map[string]string{
		"email":    email,
		"password": password,
	}, "")
	if w.Code != http.StatusOK {
		s.t.Fatalf("register failed: status=%d body=%s", w.Code, w.Body.String())
	}
	return parseAuth(s.t, w)
}

func (s *TestServer) Login(email, password string) (token string) {
	s.t.Helper()
	w := s.Do(http.MethodPost, "/api/auth/login", map[string]string{
		"email":    email,
		"password": password,
	}, "")
	if w.Code != http.StatusOK {
		s.t.Fatalf("login failed: status=%d body=%s", w.Code, w.Body.String())
	}
	tok, _, _ := parseAuth(s.t, w)
	return tok
}

func parseAuth(t *testing.T, w *httptest.ResponseRecorder) (string, uint, string) {
	t.Helper()
	var body struct {
		Token string `json:"token"`
		User  struct {
			ID   uint   `json:"id"`
			Role string `json:"role"`
		} `json:"user"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal auth response: %v", err)
	}
	return body.Token, body.User.ID, body.User.Role
}

// AsAdmin promotes the user with the given email to the admin role directly
// in the database (bypasses RBAC self-protection so tests can set up admin
// fixtures cleanly).
func (s *TestServer) AsAdmin(email string) {
	s.t.Helper()
	if err := database.DB.Exec(`UPDATE users SET role = 'admin' WHERE email = ?`, email).Error; err != nil {
		s.t.Fatalf("promote to admin: %v", err)
	}
}
