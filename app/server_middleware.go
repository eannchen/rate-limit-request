package app

import (
	"fmt"
	"net/http"
	"rate-limit-request/model"
	"rate-limit-request/response"
	"rate-limit-request/util/ginger"
	"strconv"

	"github.com/gin-gonic/gin"
)

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ginger.ExtendResponse(c, response.NewWrap())

		key := fmt.Sprintf("rete_limit_%s_%s_%s", c.Request.URL.Path, c.Request.Method, c.ClientIP())
		// RateLimit, err := CacheRepository.IncreaseRateLimit(ctx, key)
		RateLimit, err := CacheRepository.IncreaseRateLimitByLua(ctx, key)
		if err != nil {
			ctx.WithError(err).Response(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		ctx.Set("RateLimit", RateLimit)
		ctx.Header("X-RateLimit-Maximum", strconv.Itoa(model.RateLimitMaximum))
		ctx.Header("X-RateLimit-Current", strconv.Itoa(RateLimit.Count))

		if RateLimit.Count > model.RateLimitMaximum {
			ctx.WithError(err).Response(http.StatusTooManyRequests, "Too Many Requests")
			return
		}

		ctx.Next()
	}
}
