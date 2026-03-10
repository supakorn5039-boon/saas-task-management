package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/config"
	"github.com/supakorn5039-boon/saas-task-backend/src/server/routes"
)

func WebServer(config *config.ServerConfig) {
	if config.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	applyCorsMiddleware(router)
	routes.ApplyRoutes(router)

	if err := router.Run(fmt.Sprintf(":%d", config.Port)); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
