package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID injects a unique request ID and correlation ID into each request context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = requestID
		}

		c.Set("request_id", requestID)
		c.Set("correlation_id", correlationID)
		c.Header("X-Request-ID", requestID)
		c.Header("X-Correlation-ID", correlationID)

		c.Next()
	}
}

// Logger logs structured JSON for every request.
func Logger(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		slog.Info("http_request",
			"service", serviceName,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
			"request_id", c.GetString("request_id"),
			"correlation_id", c.GetString("correlation_id"),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}

// CORS adds CORS headers for cross-origin requests.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID, X-Correlation-ID, Idempotency-Key")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID, X-Correlation-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
