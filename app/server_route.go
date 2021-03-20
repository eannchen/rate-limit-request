package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type route struct {
	method      string
	endpoint    string
	middlewares []gin.HandlerFunc
	handler     gin.HandlerFunc
}

var routes = []route{
	route{http.MethodGet, "/app", []gin.HandlerFunc{rateLimitMiddleware()}, appGetHandler},
}
