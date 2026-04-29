package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"github.com/supakorn5039-boon/saas-task-backend/src/testhelpers"
)

// newTestContext returns a *gin.Context backed by an httptest recorder. Useful
// for unit-testing services that read user_id / email / IP from the context.
func newTestContext(actorID uint, actorEmail, ip, ua string) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = ip + ":12345"
	if ua != "" {
		req.Header.Set("User-Agent", ua)
	}
	c.Request = req
	if actorID != 0 {
		c.Set("user_id", actorID)
	}
	if actorEmail != "" {
		c.Set("email", actorEmail)
	}
	return c
}

func TestAuditService_Record(t *testing.T) {
	testhelpers.SetupTestDB(t)
	svc := service.NewAuditService()

	t.Run("captures actor id and email from gin context", func(t *testing.T) {
		c := newTestContext(42, "alice@example.com", "10.0.0.1", "test-agent/1.0")
		svc.Record(c, model.AuditActionLogin, model.AuditStatusSuccess, service.RecordOpts{})

		out, err := svc.ListLogs(service.ListAuditLogsOptions{Page: 1, PerPage: 10})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(out.Data) != 1 {
			t.Fatalf("want 1 row, got %d", len(out.Data))
		}
		got := out.Data[0]
		if got.ActorID == nil || *got.ActorID != 42 {
			t.Errorf("ActorID = %v, want 42", got.ActorID)
		}
		if got.ActorEmail != "alice@example.com" {
			t.Errorf("ActorEmail = %q, want alice@example.com", got.ActorEmail)
		}
		if got.Action != model.AuditActionLogin {
			t.Errorf("Action = %q, want %q", got.Action, model.AuditActionLogin)
		}
		if got.Status != model.AuditStatusSuccess {
			t.Errorf("Status = %q, want success", got.Status)
		}
		if got.UserAgent != "test-agent/1.0" {
			t.Errorf("UserAgent = %q, want test-agent/1.0", got.UserAgent)
		}
	})

	t.Run("explicit opts override context values", func(t *testing.T) {
		// A failed-login event has no authenticated actor — the controller
		// passes ActorEmail explicitly. Make sure that wins over an empty ctx.
		c := newTestContext(0, "", "10.0.0.2", "")
		svc.Record(c, model.AuditActionLoginFailed, model.AuditStatusFailure, service.RecordOpts{
			ActorEmail: "attacker@example.com",
			Metadata:   model.JSONB{"reason": "bad password"},
		})

		out, _ := svc.ListLogs(service.ListAuditLogsOptions{Page: 1, PerPage: 10, Action: model.AuditActionLoginFailed})
		if len(out.Data) != 1 {
			t.Fatalf("want 1 row, got %d", len(out.Data))
		}
		got := out.Data[0]
		if got.ActorEmail != "attacker@example.com" {
			t.Errorf("ActorEmail = %q, want attacker@example.com", got.ActorEmail)
		}
		if got.Status != model.AuditStatusFailure {
			t.Errorf("Status = %q, want failure", got.Status)
		}
		if got.Metadata["reason"] != "bad password" {
			t.Errorf("Metadata.reason = %v, want \"bad password\"", got.Metadata["reason"])
		}
	})
}

func TestAuditService_ListLogs(t *testing.T) {
	testhelpers.SetupTestDB(t)
	svc := service.NewAuditService()

	c1 := newTestContext(1, "u1@example.com", "10.0.0.1", "")
	c2 := newTestContext(2, "u2@example.com", "10.0.0.2", "")
	svc.Record(c1, model.AuditActionLogin, model.AuditStatusSuccess, service.RecordOpts{})
	svc.Record(c2, model.AuditActionLogin, model.AuditStatusSuccess, service.RecordOpts{})
	svc.Record(c1, model.AuditActionTaskCreated, model.AuditStatusSuccess, service.RecordOpts{})

	t.Run("filters by action", func(t *testing.T) {
		out, _ := svc.ListLogs(service.ListAuditLogsOptions{
			Page: 1, PerPage: 10, Action: model.AuditActionTaskCreated,
		})
		if out.Meta.Total != 1 {
			t.Errorf("Total = %d, want 1", out.Meta.Total)
		}
	})

	t.Run("filters by actor (used by my-activity endpoint)", func(t *testing.T) {
		actor := uint(1)
		out, _ := svc.ListLogs(service.ListAuditLogsOptions{
			Page: 1, PerPage: 10, ActorID: &actor,
		})
		if out.Meta.Total != 2 {
			t.Errorf("Total = %d, want 2 events for actor 1", out.Meta.Total)
		}
	})

	t.Run("filters by actor email substring", func(t *testing.T) {
		out, _ := svc.ListLogs(service.ListAuditLogsOptions{
			Page: 1, PerPage: 10, Search: "u2",
		})
		if out.Meta.Total != 1 {
			t.Errorf("Total = %d, want 1 (only u2)", out.Meta.Total)
		}
	})

	t.Run("filters by date range", func(t *testing.T) {
		future := time.Now().Add(1 * time.Hour)
		out, _ := svc.ListLogs(service.ListAuditLogsOptions{
			Page: 1, PerPage: 10, From: future,
		})
		if out.Meta.Total != 0 {
			t.Errorf("From=future should match nothing, got %d", out.Meta.Total)
		}
	})

	t.Run("paginates", func(t *testing.T) {
		out, _ := svc.ListLogs(service.ListAuditLogsOptions{Page: 1, PerPage: 2})
		if len(out.Data) != 2 {
			t.Errorf("len = %d, want 2", len(out.Data))
		}
		if out.Meta.Total != 3 {
			t.Errorf("Total = %d, want 3", out.Meta.Total)
		}
	})
}
