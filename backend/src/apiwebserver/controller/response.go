package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func successResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func errorResponse(c *gin.Context, err string, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	c.AbortWithStatusJSON(code, gin.H{"error": err})
}
