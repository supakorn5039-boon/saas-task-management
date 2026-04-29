package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
)

type AuditController struct {
	svc *service.AuditService
}

func NewAuditController() *AuditController {
	return &AuditController{svc: service.NewAuditService()}
}

func (ctrl *AuditController) RegisterRoutes(r *gin.RouterGroup) {
	// Admin view: see every event from every actor.
	admin := r.Group("/admin/audit-logs")
	admin.Use(middleware.Protected(), middleware.Rbac("admin"))
	{
		admin.GET("", ctrl.listAll)
	}

	// "My activity" view: any signed-in user can see their own audit trail.
	user := r.Group("/user/activity")
	user.Use(middleware.Protected())
	{
		user.GET("", ctrl.listMine)
	}
}

func (ctrl *AuditController) listAll(c *gin.Context) {
	opts := buildListOpts(c)
	result, err := ctrl.svc.ListLogs(opts)
	if err != nil {
		errorResponse(c, err)
		return
	}
	successResponse(c, result)
}

func (ctrl *AuditController) listMine(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	opts := buildListOpts(c)
	opts.ActorID = &userID
	result, err := ctrl.svc.ListLogs(opts)
	if err != nil {
		errorResponse(c, err)
		return
	}
	successResponse(c, result)
}

// buildListOpts turns query string params into ListAuditLogsOptions. Date
// filters use RFC3339; a malformed date is silently dropped (the filter is
// optional) — the alternative would be a confusing 400 for a bookmarked URL.
func buildListOpts(c *gin.Context) service.ListAuditLogsOptions {
	opts := service.ListAuditLogsOptions{
		Page:    parsePositiveInt(c.Query("page"), 1),
		PerPage: clampInt(parsePositiveInt(c.Query("per_page"), defaultPerPage), 1, maxPerPage),
		Action:  c.Query("action"),
		Search:  c.Query("search"),
		Sort:    c.Query("sort"),
		Order:   c.Query("order"),
	}
	if from := c.Query("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			opts.From = t
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			opts.To = t
		}
	}
	return opts
}
