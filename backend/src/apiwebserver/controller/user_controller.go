package controller

import (
	"net/http"

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
	}
}

func (ctrl *UserController) getProfile(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		errorResponse(c, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	user, err := ctrl.svc.GetUserById(userId.(uint))
	if err != nil {
		errorResponse(c, err.Error(), http.StatusNotFound)
		return
	}

	successResponse(c, user)
}
