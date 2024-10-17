package middleware

import (
	"fastgin/common/httpz"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

func RateLimitMiddleware(fillInterval time.Duration, capacity int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucket(fillInterval, capacity)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			httpz.BadRequest(c, "访问限流")
			c.Abort()
			return
		}
		c.Next()
	}
}
