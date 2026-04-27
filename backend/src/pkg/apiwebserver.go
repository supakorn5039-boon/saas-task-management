package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/controller"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/middleware"
	"github.com/supakorn5039-boon/saas-task-backend/src/config"
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
)

func MountAPIWebServer(r *gin.Engine) {
	// Initialize JWT
	security.InitJWT(config.App.Server.JWTSecret)

	// Structured request logging with request id (replaces gin's default logger).
	r.Use(middleware.RequestLogger())

	// Conservative security headers on every response (XCTO, XFO, etc).
	r.Use(middleware.SecurityHeaders())

	// Apply CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://saas-management.local"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", middleware.RequestIDHeader},
		ExposeHeaders:    []string{"Content-Length", middleware.RequestIDHeader},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Mount controllers
	api := r.Group("/api")

	authCtrl := controller.NewAuthController()
	authCtrl.RegisterRoutes(api)

	userCtrl := controller.NewUserController()
	userCtrl.RegisterRoutes(api)

	taskCtrl := controller.NewTaskController()
	taskCtrl.RegisterRoutes(api)

	adminCtrl := controller.NewAdminController()
	adminCtrl.RegisterRoutes(api)

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// /healthz is a real readiness probe — we ping the DB so a load balancer
	// pulls the pod out of rotation if Postgres is unreachable. /ping above
	// stays as a lighter liveness check.
	api.GET("/healthz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := database.DB.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "db": "unreachable"})
			return
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "down", "db": "ping_failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "ok"})
	})
}
