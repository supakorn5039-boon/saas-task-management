package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONB is a thin wrapper that lets us store/load arbitrary JSON in a Postgres
// JSONB column via GORM. Marshals to []byte on write, unmarshals on read.
type JSONB map[string]any

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(src any) error {
	if src == nil {
		*j = nil
		return nil
	}
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return errors.New("audit metadata: unsupported scan source")
	}
	return json.Unmarshal(b, j)
}

// AuditLog is append-only. No gorm.Model — we don't want updated_at or
// soft-delete. The actor email is denormalized so the log stays readable
// after a user is deleted.
type AuditLog struct {
	ID          uint      `gorm:"primaryKey"`
	ActorUserID *uint     `gorm:"index"`
	ActorEmail  string    `gorm:"size:255"`
	Action      string    `gorm:"size:64;not null;index"`
	TargetType  string    `gorm:"size:32"`
	TargetID    *uint     ``
	Status      string    `gorm:"size:16;not null;default:success"`
	IP          string    `gorm:"size:64"`
	UserAgent   string    `gorm:"size:512"`
	Metadata    JSONB     `gorm:"type:jsonb"`
	CreatedAt   time.Time `gorm:"not null"`
}

func (AuditLog) TableName() string { return "audit_logs" }

// Audit status — kept as exported constants so callers don't pass typo'd
// strings. "success" is the default; "failure" marks attempts that didn't
// complete (e.g. login with wrong password).
const (
	AuditStatusSuccess = "success"
	AuditStatusFailure = "failure"
)

// Action identifiers — namespaced as "<domain>.<verb>" to keep filtering
// straightforward in the admin UI.
const (
	AuditActionLogin           = "auth.login"
	AuditActionLoginFailed     = "auth.login_failed"
	AuditActionRegister        = "auth.register"
	AuditActionLogout          = "auth.logout"
	AuditActionPasswordChanged = "user.password_changed"
	AuditActionUserUpdated     = "admin.user_updated"
	AuditActionUserDeleted     = "admin.user_deleted"
	AuditActionTaskCreated     = "task.created"
	AuditActionTaskUpdated     = "task.updated"
	AuditActionTaskDeleted     = "task.deleted"
)

type AuditLogDto struct {
	ID         uint           `json:"id"`
	ActorID    *uint          `json:"actorId,omitempty"`
	ActorEmail string         `json:"actorEmail"`
	Action     string         `json:"action"`
	TargetType string         `json:"targetType,omitempty"`
	TargetID   *uint          `json:"targetId,omitempty"`
	Status     string         `json:"status"`
	IP         string         `json:"ip,omitempty"`
	UserAgent  string         `json:"userAgent,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
	CreatedAt  string         `json:"createdAt"`
}

func (a *AuditLog) ToDto() *AuditLogDto {
	var meta map[string]any
	if a.Metadata != nil {
		meta = a.Metadata
	}
	return &AuditLogDto{
		ID:         a.ID,
		ActorID:    a.ActorUserID,
		ActorEmail: a.ActorEmail,
		Action:     a.Action,
		TargetType: a.TargetType,
		TargetID:   a.TargetID,
		Status:     a.Status,
		IP:         a.IP,
		UserAgent:  a.UserAgent,
		Metadata:   meta,
		CreatedAt:  a.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type AuditLogListMeta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"perPage"`
	Total   int64 `json:"total"`
}

type AuditLogListResponse struct {
	Data []*AuditLogDto   `json:"data"`
	Meta AuditLogListMeta `json:"meta"`
}
