package integration

import (
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestAdmin_RBAC_BlocksRegularUser(t *testing.T) {
	srv := newServer(t)
	userToken, _, _ := srv.Register("user@example.com", "password123")

	// Regular user — must be blocked from every admin endpoint.
	cases := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/admin/users"},
		{http.MethodPut, "/api/admin/users/1"},
		{http.MethodDelete, "/api/admin/users/1"},
	}
	for _, c := range cases {
		t.Run(c.method+" "+c.path, func(t *testing.T) {
			w := srv.Do(c.method, c.path, nil, userToken)
			if w.Code != http.StatusForbidden {
				t.Errorf("status = %d, want 403", w.Code)
			}
		})
	}
}

func TestAdmin_ListUsers(t *testing.T) {
	srv := newServer(t)
	srv.Register("alpha@example.com", "password123")
	srv.Register("beta@example.com", "password123")
	_, _, _ = srv.Register("admin@example.com", "password123")
	srv.AsAdmin("admin@example.com")
	adminToken := srv.Login("admin@example.com", "password123")

	w := srv.Do(http.MethodGet, "/api/admin/users", nil, adminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}

	var body struct {
		Data []map[string]any `json:"data"`
		Meta struct{ Total int }
	}
	mustUnmarshal(t, w.Body.Bytes(), &body)
	if body.Meta.Total != 3 {
		t.Errorf("total = %d, want 3", body.Meta.Total)
	}
}

func TestAdmin_UpdateUser_PromoteAndProtect(t *testing.T) {
	srv := newServer(t)
	srv.Register("admin@example.com", "password123")
	srv.AsAdmin("admin@example.com")
	adminToken := srv.Login("admin@example.com", "password123")

	_, memberID, _ := srv.Register("member@example.com", "password123")

	t.Run("promote member to manager", func(t *testing.T) {
		w := srv.Do(http.MethodPut, fmt.Sprintf("/api/admin/users/%d", memberID), map[string]any{
			"role": "manager",
		}, adminToken)
		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
		}
	})

	// adminID has to be looked up via the users list (id=1 isn't guaranteed
	// because tests use truncate, but the registration order is deterministic
	// within this test).
	w := srv.Do(http.MethodGet, "/api/admin/users", nil, adminToken)
	var body struct {
		Data []map[string]any `json:"data"`
	}
	mustUnmarshal(t, w.Body.Bytes(), &body)
	var adminID int
	for _, u := range body.Data {
		if u["email"] == "admin@example.com" {
			adminID = int(u["id"].(float64))
		}
	}
	if adminID == 0 {
		t.Fatal("could not locate admin user in list")
	}

	t.Run("admin cannot demote themselves out of admin", func(t *testing.T) {
		w := srv.Do(http.MethodPut, "/api/admin/users/"+strconv.Itoa(adminID), map[string]any{
			"role": "user",
		}, adminToken)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("admin cannot deactivate themselves", func(t *testing.T) {
		w := srv.Do(http.MethodPut, "/api/admin/users/"+strconv.Itoa(adminID), map[string]any{
			"status": 0,
		}, adminToken)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("admin cannot delete themselves", func(t *testing.T) {
		w := srv.Do(http.MethodDelete, "/api/admin/users/"+strconv.Itoa(adminID), nil, adminToken)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})

	t.Run("invalid role rejected with 400", func(t *testing.T) {
		w := srv.Do(http.MethodPut, fmt.Sprintf("/api/admin/users/%d", memberID), map[string]any{
			"role": "superuser",
		}, adminToken)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})
}

func TestAdmin_DeleteMember(t *testing.T) {
	srv := newServer(t)
	srv.Register("admin@example.com", "password123")
	srv.AsAdmin("admin@example.com")
	adminToken := srv.Login("admin@example.com", "password123")

	_, victimID, _ := srv.Register("victim@example.com", "password123")

	w := srv.Do(http.MethodDelete, fmt.Sprintf("/api/admin/users/%d", victimID), nil, adminToken)
	if w.Code != http.StatusOK {
		t.Fatalf("delete: status = %d, body = %s", w.Code, w.Body.String())
	}

	// Second delete → 404 (already soft-deleted).
	w = srv.Do(http.MethodDelete, fmt.Sprintf("/api/admin/users/%d", victimID), nil, adminToken)
	if w.Code != http.StatusNotFound {
		t.Errorf("second delete: status = %d, want 404", w.Code)
	}
}
