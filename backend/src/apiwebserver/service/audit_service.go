// Package service — audit_service handles the append-only audit log:
// recording events (used by other controllers) and querying them (admin
// view + per-user "my activity" view). Reads/writes go through the same
// gorm.DB used by the rest of the services.
package service

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apperror"
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"gorm.io/gorm"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService() *AuditService {
	return &AuditService{database.DB}
}

// RecordOpts is everything optional about an audit entry. Action and Status
// are required and live on the top-level Record signature.
type RecordOpts struct {
	ActorID    *uint
	ActorEmail string
	TargetType string
	TargetID   *uint
	Metadata   model.JSONB
}

// Record writes an audit entry. Failures are logged and swallowed — auditing
// must never break the user-facing request. The actor (id + email) and IP/UA
// are pulled from the gin.Context when available so callers don't have to
// thread them through every layer.
func (s *AuditService) Record(c *gin.Context, action, status string, opts RecordOpts) {
	entry := model.AuditLog{
		Action:     action,
		Status:     status,
		TargetType: opts.TargetType,
		TargetID:   opts.TargetID,
		Metadata:   opts.Metadata,
		CreatedAt:  time.Now(),
	}

	entry.ActorUserID = opts.ActorID
	entry.ActorEmail = opts.ActorEmail

	if c != nil {
		if entry.ActorUserID == nil {
			if v, ok := c.Get("user_id"); ok {
				if uid, ok := v.(uint); ok {
					entry.ActorUserID = &uid
				}
			}
		}
		// Fall back to the JWT email — Protected middleware sets it on every
		// authenticated request, so admin actions don't need to load the user.
		if entry.ActorEmail == "" {
			if v, ok := c.Get("email"); ok {
				if em, ok := v.(string); ok {
					entry.ActorEmail = em
				}
			}
		}
		entry.IP = c.ClientIP()
		entry.UserAgent = truncate(c.Request.UserAgent(), 512)
	}

	if err := s.db.Create(&entry).Error; err != nil {
		slog.Error("audit_log_write_failed",
			"action", action,
			"status", status,
			"err", err.Error(),
		)
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

// ----- Listing (admin + my activity) -----

// allowedAuditSortColumns intentionally limits sort options — keeps query
// shape predictable and avoids leaking column names from arbitrary input.
var allowedAuditSortColumns = map[string]string{
	"created_at": "created_at",
	"action":     "action",
	"status":     "status",
}

type ListAuditLogsOptions struct {
	Page    int
	PerPage int
	// Filters — all optional. Empty string / zero value disables the filter.
	Action string
	Search string // matches actor_email ILIKE
	From   time.Time
	To     time.Time
	// Sort
	Sort  string
	Order string
	// Restrict to a single actor (used by the "my activity" endpoint).
	ActorID *uint
}

func (s *AuditService) ListLogs(opts ListAuditLogsOptions) (*model.AuditLogListResponse, error) {
	q := s.db.Model(&model.AuditLog{})

	if opts.ActorID != nil {
		q = q.Where("actor_user_id = ?", *opts.ActorID)
	}
	if opts.Action != "" {
		q = q.Where("action = ?", opts.Action)
	}
	if opts.Search != "" {
		q = q.Where("actor_email ILIKE ?", "%"+opts.Search+"%")
	}
	if !opts.From.IsZero() {
		q = q.Where("created_at >= ?", opts.From)
	}
	if !opts.To.IsZero() {
		q = q.Where("created_at <= ?", opts.To)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, apperror.Wrap(err, 500, "failed to count audit logs")
	}

	sortCol, ok := allowedAuditSortColumns[opts.Sort]
	if !ok {
		sortCol = "created_at"
	}
	order := "desc"
	if opts.Order == "asc" {
		order = "asc"
	}

	var rows []model.AuditLog
	err := q.Order(sortCol + " " + order).
		Limit(opts.PerPage).
		Offset((opts.Page - 1) * opts.PerPage).
		Find(&rows).Error
	if err != nil {
		return nil, apperror.Wrap(err, 500, "failed to list audit logs")
	}

	dtos := make([]*model.AuditLogDto, len(rows))
	for i, r := range rows {
		dtos[i] = r.ToDto()
	}

	return &model.AuditLogListResponse{
		Data: dtos,
		Meta: model.AuditLogListMeta{
			Page:    opts.Page,
			PerPage: opts.PerPage,
			Total:   total,
		},
	}, nil
}
