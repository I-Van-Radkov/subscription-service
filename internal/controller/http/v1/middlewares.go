package v1

import (
	"github.com/I-Van-Radkov/subscription-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.Request.Header.Get("x-request-id")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		ctx := logger.WithRequestID(c.Request.Context(), reqID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
