package app

import (
	"fmt"
	"net/http"
	"rate-limit-request/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("rete_limit_%s_%s_%s", c.Request.URL.Path, c.Request.Method, c.ClientIP())
		// RateLimit, err := CacheRepository.IncreaseRateLimit(c, key)
		RateLimit, err := CacheRepository.IncreaseRateLimitByLua(c, key)
		if err != nil {
			c.String(http.StatusInternalServerError, "Error")
			c.Abort()
			return
		}

		c.Set("RateLimit", RateLimit)
		c.Header("X-RateLimit-Maximum", strconv.Itoa(model.RateLimitMaximum))
		c.Header("X-RateLimit-Current", strconv.Itoa(RateLimit.Count))

		if RateLimit.Count > model.RateLimitMaximum {
			c.String(http.StatusTooManyRequests, "Error")
			c.Abort()
			return
		}

		c.Next()
	}
}
