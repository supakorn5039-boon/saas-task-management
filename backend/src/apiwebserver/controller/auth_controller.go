package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
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
	// 5 attempts / IP / minute — enough for a real user, slow enough to make
	// credential stuffing unattractive.
	auth.Use(middleware.RateLimit(5, time.Minute))
	{
		auth.POST("/login", ctrl.login)
		auth.POST("/register", ctrl.register)
	}
}

// loginCredentials uses min=1 so existing accounts with short passwords still
// work; min=8 is enforced at registration time only.
type loginCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type registerCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func (ctrl *AuthController) login(c *gin.Context) {
	var body loginCredentials
	if err := c.ShouldBindJSON(&body); err != nil {
		badRequest(c, err.Error())
		return
	}

	user, err := ctrl.svc.Login(body.Email, body.Password)
	if err != nil {
		errorResponse(c, err)
		return
	}

	issueToken(c, user)
}

func (ctrl *AuthController) register(c *gin.Context) {
	var body registerCredentials
	if err := c.ShouldBindJSON(&body); err != nil {
		badRequest(c, err.Error())
		return
	}

	user, err := ctrl.svc.Register(body.Email, body.Password)
	if err != nil {
		errorResponse(c, err)
		return
	}

	issueToken(c, user)
}

func issueToken(c *gin.Context, user *model.UserDto) {
	token, err := security.GenerateJWT(user.Id, user.Role)
	if err != nil {
		errorResponse(c, err)
		return
	}
	successResponse(c, gin.H{"token": token, "user": user})
}
