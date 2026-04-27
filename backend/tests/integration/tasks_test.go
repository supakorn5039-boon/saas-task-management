package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"testing"
)

func TestTasks_RequiresAuth(t *testing.T) {
	srv := newServer(t)

	w := srv.Do(http.MethodGet, "/api/tasks", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("unauthenticated GET /api/tasks: status = %d, want 401", w.Code)
	}
}

func TestTasks_CRUD(t *testing.T) {
	srv := newServer(t)
	token, _, _ := srv.Register("owner@example.com", "password123")

	// Create
	w := srv.Do(http.MethodPost, "/api/tasks", map[string]any{
		"title":       "Buy milk",
		"description": "skim",
	}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("create: status = %d, body = %s", w.Code, w.Body.String())
	}
	id := decodeID(t, w.Body.Bytes())

	// List
	w = srv.Do(http.MethodGet, "/api/tasks", nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("list: status = %d", w.Code)
	}
	var listBody struct {
		Data   []map[string]any `json:"data"`
		Counts map[string]int   `json:"counts"`
	}
	mustUnmarshal(t, w.Body.Bytes(), &listBody)
	if len(listBody.Data) != 1 {
		t.Fatalf("expected 1 task, got %d", len(listBody.Data))
	}
	if listBody.Counts["all"] != 1 || listBody.Counts["todo"] != 1 {
		t.Errorf("counts = %+v, want all=1 todo=1", listBody.Counts)
	}

	// Patch (full edit + status change)
	w = srv.Do(http.MethodPut, "/api/tasks/"+strconv.Itoa(id), map[string]any{
		"title":  "Buy milk (edited)",
		"status": "in_progress",
	}, token)
	if w.Code != http.StatusOK {
		t.Fatalf("update: status = %d, body = %s", w.Code, w.Body.String())
	}

	// Empty patch is rejected (defends against client bugs sending {}).
	w = srv.Do(http.MethodPut, "/api/tasks/"+strconv.Itoa(id), map[string]any{}, token)
	if w.Code != http.StatusBadRequest {
		t.Errorf("empty update: status = %d, want 400", w.Code)
	}

	// Delete
	w = srv.Do(http.MethodDelete, "/api/tasks/"+strconv.Itoa(id), nil, token)
	if w.Code != http.StatusOK {
		t.Fatalf("delete: status = %d", w.Code)
	}

	// Second delete → 404
	w = srv.Do(http.MethodDelete, "/api/tasks/"+strconv.Itoa(id), nil, token)
	if w.Code != http.StatusNotFound {
		t.Errorf("second delete: status = %d, want 404", w.Code)
	}
}

func TestTasks_TenantIsolation(t *testing.T) {
	// Two users; alice's task must not be visible / mutable by bob.
	srv := newServer(t)
	aliceToken, _, _ := srv.Register("alice@example.com", "password123")
	bobToken, _, _ := srv.Register("bob@example.com", "password123")

	w := srv.Do(http.MethodPost, "/api/tasks", map[string]any{"title": "Alice's task"}, aliceToken)
	if w.Code != http.StatusOK {
		t.Fatalf("alice create: %d", w.Code)
	}
	id := decodeID(t, w.Body.Bytes())

	t.Run("bob cannot see alice's task in his list", func(t *testing.T) {
		w := srv.Do(http.MethodGet, "/api/tasks", nil, bobToken)
		var body struct {
			Data []map[string]any `json:"data"`
		}
		mustUnmarshal(t, w.Body.Bytes(), &body)
		if len(body.Data) != 0 {
			t.Errorf("bob sees %d tasks, want 0", len(body.Data))
		}
	})

	t.Run("bob cannot update alice's task — gets 404 (not 403)", func(t *testing.T) {
		// 404 over 403 is intentional: 403 would leak that the task exists.
		w := srv.Do(http.MethodPut, fmt.Sprintf("/api/tasks/%d", id), map[string]any{
			"status": "done",
		}, bobToken)
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})

	t.Run("bob cannot delete alice's task — also 404", func(t *testing.T) {
		w := srv.Do(http.MethodDelete, fmt.Sprintf("/api/tasks/%d", id), nil, bobToken)
		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want 404", w.Code)
		}
	})
}

func TestTasks_PaginationFilterSearch(t *testing.T) {
	srv := newServer(t)
	token, _, _ := srv.Register("seeder@example.com", "password123")

	// Seed 12 tasks across statuses.
	for i := 1; i <= 12; i++ {
		body := map[string]any{
			"title":       fmt.Sprintf("Task %d", i),
			"description": fmt.Sprintf("desc %d", i),
		}
		w := srv.Do(http.MethodPost, "/api/tasks", body, token)
		if w.Code != http.StatusOK {
			t.Fatalf("seed task %d: %d", i, w.Code)
		}
	}

	t.Run("pagination respects per_page", func(t *testing.T) {
		w := srv.Do(http.MethodGet, "/api/tasks?per_page=5", nil, token)
		var body struct {
			Data []map[string]any `json:"data"`
			Meta struct{ Total int }
		}
		mustUnmarshal(t, w.Body.Bytes(), &body)
		if len(body.Data) != 5 {
			t.Errorf("len(data) = %d, want 5", len(body.Data))
		}
		if body.Meta.Total != 12 {
			t.Errorf("total = %d, want 12", body.Meta.Total)
		}
	})

	t.Run("search is case-insensitive", func(t *testing.T) {
		w := srv.Do(http.MethodGet, "/api/tasks?search=TASK%201", nil, token)
		var body struct {
			Data []map[string]any `json:"data"`
		}
		mustUnmarshal(t, w.Body.Bytes(), &body)
		// "Task 1" matches "Task 1", "Task 10", "Task 11", "Task 12" → 4 results
		if len(body.Data) < 1 {
			t.Errorf("search: got %d results, want >= 1", len(body.Data))
		}
	})

	t.Run("invalid status filter rejected with 400", func(t *testing.T) {
		w := srv.Do(http.MethodGet, "/api/tasks?status=banana", nil, token)
		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want 400", w.Code)
		}
	})
}

// ----- helpers -----

func decodeID(t *testing.T, body []byte) int {
	t.Helper()
	var v struct {
		ID int `json:"id"`
	}
	mustUnmarshal(t, body, &v)
	if v.ID == 0 {
		t.Fatalf("expected id in body, got: %s", body)
	}
	return v.ID
}

func mustUnmarshal(t *testing.T, body []byte, into any) {
	t.Helper()
	if err := json.Unmarshal(body, into); err != nil {
		t.Fatalf("unmarshal: %v body=%s", err, body)
	}
}
