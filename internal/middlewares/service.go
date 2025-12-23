package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
)

func SetupMiddlewares(r *gin.Engine) {
	r.Use(minDelay(200 * time.Millisecond))
}
