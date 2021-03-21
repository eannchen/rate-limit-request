package app

import (
	"net/http"
	"rate-limit-request/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func appGetHandler(c *gin.Context) {
	data := c.MustGet("RateLimit").(model.RateLimit)
	c.String(http.StatusOK, strconv.Itoa(data.Count))
	return
}
