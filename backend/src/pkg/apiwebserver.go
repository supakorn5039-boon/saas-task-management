package pkg

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apiwebserver/controller"
	"github.com/supakorn5039-boon/saas-task-backend/src/config"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
)

func MountAPIWebServer(r *gin.Engine) {
	// Initialize JWT
	security.InitJWT(config.App.Server.JWTSecret)

	// Apply CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
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

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
}
