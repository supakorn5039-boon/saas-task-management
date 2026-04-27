package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supakorn5039-boon/saas-task-backend/src/apperror"
)

func successResponse(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// errorResponse maps an error to a JSON response. AppError carries its own
// status + client-safe message; anything else becomes a generic 500 and the
// real error is logged (never sent to the client).
func errorResponse(c *gin.Context, err error) {
	if ae, ok := apperror.As(err); ok {
		if ae.Err != nil {
			log.Printf("request error: status=%d msg=%q err=%v", ae.Status, ae.Message, ae.Err)
		}
		c.AbortWithStatusJSON(ae.Status, gin.H{"error": ae.Message})
		return
	}

	log.Printf("unhandled error: %v", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

// badRequest is a shortcut for 400s with a known message (e.g. validation).
func badRequest(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": message})
}
