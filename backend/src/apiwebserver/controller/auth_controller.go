package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/service"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
)

type AuthController struct {
	svc *service.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{svc: service.NewAuthService()}
}

func (ctrl *AuthController) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", ctrl.login)
		auth.POST("/register", ctrl.register)
	}
}

func (ctrl *AuthController) login(c *gin.Context) {
	var body model.CredentialDto
	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ctrl.svc.Login(body.Email, body.Password)
	if err != nil {
		errorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := security.GenerateJWT(user.Id)
	if err != nil {
		errorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	successResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func (ctrl *AuthController) register(c *gin.Context) {
	var body model.CredentialDto
	if err := c.ShouldBindJSON(&body); err != nil {
		errorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := ctrl.svc.Register(body.Email, body.Password)
	if err != nil {
		errorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := security.GenerateJWT(user.Id)
	if err != nil {
		errorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	successResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}
