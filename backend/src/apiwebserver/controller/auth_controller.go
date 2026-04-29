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
	svc   *service.AuthService
	audit *service.AuditService
}

func NewAuthController() *AuthController {
	return &AuthController{
		svc:   service.NewAuthService(),
		audit: service.NewAuditService(),
	}
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

	// Logout requires a valid token — the endpoint is mostly bookkeeping
	// (JWT is stateless) but it lets us record the event in the audit log.
	logout := r.Group("/auth/logout")
	logout.Use(middleware.Protected())
	{
		logout.POST("", ctrl.logout)
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
		// Record the failure with the attempted email so brute-force shows up
		// as a cluster on the audit page. No actor id (login wasn't valid).
		ctrl.audit.Record(c, model.AuditActionLoginFailed, model.AuditStatusFailure, service.RecordOpts{
			ActorEmail: body.Email,
			Metadata:   model.JSONB{"reason": err.Error()},
		})
		errorResponse(c, err)
		return
	}

	ctrl.audit.Record(c, model.AuditActionLogin, model.AuditStatusSuccess, service.RecordOpts{
		ActorID:    &user.Id,
		ActorEmail: user.Email,
	})
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

	ctrl.audit.Record(c, model.AuditActionRegister, model.AuditStatusSuccess, service.RecordOpts{
		ActorID:    &user.Id,
		ActorEmail: user.Email,
	})
	issueToken(c, user)
}

func (ctrl *AuthController) logout(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	ctrl.audit.Record(c, model.AuditActionLogout, model.AuditStatusSuccess, service.RecordOpts{
		ActorID: &userID,
	})
	successResponse(c, gin.H{"message": "Logged out"})
}

func issueToken(c *gin.Context, user *model.UserDto) {
	token, err := security.GenerateJWT(user.Id, user.Role, user.Email)
	if err != nil {
		errorResponse(c, err)
		return
	}
	successResponse(c, gin.H{"token": token, "user": user})
}
