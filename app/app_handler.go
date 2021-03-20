package app

import (
	"net/http"
	"rate-limit-request/response"
	"rate-limit-request/util/ginger"

	"github.com/gin-gonic/gin"
)

func appGetHandler(c *gin.Context) {
	ctx := ginger.ExtendResponse(c, response.NewWrap())

	data := ctx.MustGet("RateLimit")
	ctx.WithData(data).Response(http.StatusOK, "")
	return
}
