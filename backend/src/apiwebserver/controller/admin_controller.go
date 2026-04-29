package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
)

type AdminController struct {
	users *service.UserService
	audit *service.AuditService
}

func NewAdminController() *AdminController {
	return &AdminController{
		users: service.NewUserService(),
		audit: service.NewAuditService(),
	}
}

func (ctrl *AdminController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	admin.Use(middleware.Protected(), middleware.Rbac("admin"))
	{
		admin.GET("/users", ctrl.listUsers)
		admin.PUT("/users/:id", ctrl.updateUser)
		admin.DELETE("/users/:id", ctrl.deleteUser)
	}
}

func (ctrl *AdminController) listUsers(c *gin.Context) {
	opts := service.ListUsersOptions{
		Page:    parsePositiveInt(c.Query("page"), 1),
		PerPage: clampInt(parsePositiveInt(c.Query("per_page"), defaultPerPage), 1, maxPerPage),
		Search:  c.Query("search"),
		Sort:    c.Query("sort"),
		Order:   c.Query("order"),
	}
	result, err := ctrl.users.ListUsers(opts)
	if err != nil {
		errorResponse(c, err)
		return
	}
	successResponse(c, result)
}

func (ctrl *AdminController) updateUser(c *gin.Context) {
	actorID := c.MustGet("user_id").(uint)
	targetID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var body struct {
		Role   *string `json:"role"`
		Status *int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		badRequest(c, err.Error())
		return
	}

	user, err := ctrl.users.AdminUpdateUser(actorID, targetID, service.AdminUpdateUserInput{
		Role:   body.Role,
		Status: body.Status,
	})
	if err != nil {
		errorResponse(c, err)
		return
	}

	meta := model.JSONB{"targetEmail": user.Email}
	if body.Role != nil {
		meta["role"] = *body.Role
	}
	if body.Status != nil {
		meta["status"] = *body.Status
	}
	ctrl.audit.Record(c, model.AuditActionUserUpdated, model.AuditStatusSuccess, service.RecordOpts{
		TargetType: "user",
		TargetID:   &user.Id,
		Metadata:   meta,
	})
	successResponse(c, user)
}

func (ctrl *AdminController) deleteUser(c *gin.Context) {
	actorID := c.MustGet("user_id").(uint)
	targetID, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	// Snapshot the target's email before delete so the audit row can show it
	// after the user row is gone (soft-deleted, but not loaded by default).
	var targetEmail string
	if t, terr := ctrl.users.GetUserById(targetID); terr == nil {
		targetEmail = t.Email
	}

	if err := ctrl.users.AdminDeleteUser(actorID, targetID); err != nil {
		errorResponse(c, err)
		return
	}

	ctrl.audit.Record(c, model.AuditActionUserDeleted, model.AuditStatusSuccess, service.RecordOpts{
		TargetType: "user",
		TargetID:   &targetID,
		Metadata:   model.JSONB{"targetEmail": targetEmail},
	})
	successResponse(c, gin.H{"message": "User deleted"})
}

// parseUintParam pulls a uint route param and writes a 400 if it's malformed.
// Returns ok=false so the caller can short-circuit.
func parseUintParam(c *gin.Context, key string) (uint, error) {
	raw := c.Param(key)
	n, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		badRequest(c, "invalid "+key)
		return 0, err
	}
	return uint(n), nil
}
