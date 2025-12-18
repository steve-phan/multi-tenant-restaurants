package middleware

import (
	"time"

	"restaurant-backend/internal/logger"

	"github.com/gin-gonic/gin"
)

// RequestLogger returns a middleware that logs HTTP requests using the application logger
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate metrics
		param := gin.LogFormatterParams{
			Request: c.Request,
			Keys:    c.Keys,
		}

		// Update logged path with query params if present
		if raw != "" {
			path = path + "?" + raw
		}

		// Log using our zap logger wrapper
		// StatusCode check to determine log level could be improved here if needed
		// For now we use the standardized LogRequest function
		logger.LogRequest(
			param.Request.Method,
			path,
			c.Writer.Status(),
			time.Since(start),
		)
	}
}
