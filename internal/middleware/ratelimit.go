package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"rate-limiter/internal/limiter"
)

func TokenBucketLimit(l *limiter.TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.ClientIP()

		result := l.Allow(user)

		// headers
		c.Header("X-RateLimit-Limit", "100")
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", result.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", result.Reset))

		if !result.Allowed {
			retryAfter := result.Reset - time.Now().Unix()

			c.Header("Retry-After", fmt.Sprintf("%d", retryAfter))
			c.JSON(429, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}