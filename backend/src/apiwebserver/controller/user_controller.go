package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
)

type UserController struct {
	svc *service.UserService
}

func NewUserController() *UserController {
	return &UserController{svc: service.NewUserService()}
}

func (ctrl *UserController) RegisterRoutes(r *gin.RouterGroup) {
	user := r.Group("/user")
	user.Use(middleware.Protected())
	{
		user.GET("/profile", ctrl.getProfile)
		user.PUT("/password", ctrl.changePassword)
	}
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

	successResponse(c, gin.H{"message": "Password updated"})
}
