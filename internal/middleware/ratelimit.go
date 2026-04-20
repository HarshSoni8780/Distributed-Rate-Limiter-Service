package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"rate-limiter/internal/limiter"
)

func RateLimimt(l *limiter.FixedWindow) gin.HandlerFunc{
	return func(c *gin.Context){
		user := c.ClientIP()
		
		if !l.Allow(user){
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}