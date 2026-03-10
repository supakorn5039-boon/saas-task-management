package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/server/controllers"
	"github.com/supakorn5039-boon/saas-task-backend/src/server/middleware"
)

func ApplyRoutes(r *gin.Engine) {

	api := r.Group("/api")
	{
		api.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.Register)
		}

		user := api.Group("/user")
		user.Use(middleware.Protected())
		{
			user.GET("/profile", controllers.GetProfile)
		}

		tasks := api.Group("/tasks")
		tasks.Use(middleware.Protected())
		{
			tasks.GET("", controllers.GetTasks)
			tasks.POST("", controllers.CreateTask)
			tasks.PATCH("/:id", controllers.ToggleTask)
			tasks.DELETE("/:id", controllers.DeleteTask)
		}
	}
}
