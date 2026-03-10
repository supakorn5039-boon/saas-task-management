package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/models"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
	"github.com/supakorn5039-boon/saas-task-backend/src/server/services"
	"github.com/supakorn5039-boon/saas-task-backend/src/utils"
)

func Login(c *gin.Context) {
	var body models.CredentialDto
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	authService := services.NewAuthenticationService()
	user, err := authService.Login(body.Email, body.Password)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := security.GenerateJWT(user.Id)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func Register(c *gin.Context) {
	var body models.CredentialDto
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	authService := services.NewAuthenticationService()
	user, err := authService.Register(body.Email, body.Password)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := security.GenerateJWT(user.Id)
	if err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"token": token,
		"user":  user,
	})
}
