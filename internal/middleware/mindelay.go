package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func MinDelay(minDuration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		elapsed := time.Since(start)
		if elapsed < minDuration {
			time.Sleep(minDuration - elapsed)
		}
	}
}
