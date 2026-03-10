package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/server/services"
	"github.com/supakorn5039-boon/saas-task-backend/src/utils"
)

func GetProfile(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, "user_id not found in context", http.StatusUnauthorized)
		return
	}

	userService := services.NewUserService()
	user, err := userService.GetUserById(userId.(uint))
	if err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusNotFound)
		return
	}

	utils.SuccessResponse(c, user)
}
