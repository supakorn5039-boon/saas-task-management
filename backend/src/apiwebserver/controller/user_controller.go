package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
)

type UserController struct {
	svc   *service.UserService
	audit *service.AuditService
}

func NewUserController() *UserController {
	return &UserController{
		svc:   service.NewUserService(),
		audit: service.NewAuditService(),
	}
}

func (ctrl *UserController) RegisterRoutes(r *gin.RouterGroup) {
	user := r.Group("/user")
	user.Use(middleware.Protected())
	{
		user.GET("/profile", ctrl.getProfile)
		user.PUT("/password", ctrl.changePassword)
	}

	// Lightweight user list for the assignee dropdown — any authenticated
	// user can read it (no PII beyond email is exposed by UserDto).
	users := r.Group("/users")
	users.Use(middleware.Protected())
	{
		users.GET("/assignable", ctrl.listAssignable)
	}
}

func (ctrl *UserController) listAssignable(c *gin.Context) {
	users, err := ctrl.svc.ListAssignable()
	if err != nil {
		errorResponse(c, err)
		return
	}
	successResponse(c, gin.H{"data": users})
}

func (ctrl *UserController) getProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	user, err := ctrl.svc.GetUserById(userID)
	if err != nil {
		errorResponse(c, err)
		return
	}

	successResponse(c, user)
}

func (ctrl *UserController) changePassword(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var body struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		badRequest(c, err.Error())
		return
	}
	if body.CurrentPassword == body.NewPassword {
		badRequest(c, "new password must be different from current")
		return
	}

	if err := ctrl.svc.ChangePassword(userID, body.CurrentPassword, body.NewPassword); err != nil {
		errorResponse(c, err)
		return
	}

	// Pull the actor's email so the audit row remains readable when the user
	// row is later deleted. Cheap — service has loaded the user already.
	if profile, perr := ctrl.svc.GetUserById(userID); perr == nil {
		ctrl.audit.Record(c, model.AuditActionPasswordChanged, model.AuditStatusSuccess, service.RecordOpts{
			ActorID:    &userID,
			ActorEmail: profile.Email,
		})
	}

	successResponse(c, gin.H{"message": "Password updated"})
}
