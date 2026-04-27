package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

const RequestIDHeader = "X-Request-ID"

// RequestLogger logs every request with a stable request id, plus the user id
// (when Protected has run). Replaces the default Gin logger so each line is
// structured and parseable.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(RequestIDHeader)
		if reqID == "" {
			reqID = newRequestID()
		}
		c.Set("request_id", reqID)
		c.Writer.Header().Set(RequestIDHeader, reqID)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		userID, _ := c.Get("user_id")
		slog.Info("http",
			"request_id", reqID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"user_id", userID,
			"client_ip", c.ClientIP(),
		)
	}
}

func newRequestID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "req-unknown"
	}
	return hex.EncodeToString(b)
}
