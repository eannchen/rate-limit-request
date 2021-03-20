package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"rate-limit-request/config"
	"rate-limit-request/storage"
	"rate-limit-request/storage/redis"
	"rate-limit-request/util/logger"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// Config config
	Config config.Config
	// CacheRepository cache repository
	CacheRepository storage.CacheRepository
)

func init() {
	var err error
	if Config, err = config.LoadConfig(); err != nil {
		panic(err)
	}

	ctx := context.Background()
	CacheRepository = redis.NewCacheRepository(Config)
	if err := CacheRepository.Init(ctx); err != nil {
		panic(err)
	}
	if err := CacheRepository.FlushDB(ctx); err != nil {
		panic(err)
	}
	logger.SetLevel(Config.Core.Mode)
}

// Run run server
func Run() {
	engine := RouterEngine()

	srv := &http.Server{
		Addr:           ":" + Config.Core.Port,
		Handler:        engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	waitShutDown(srv)

	return
}

func RouterEngine() *gin.Engine {
	gin.SetMode(Config.Core.Mode)

	router := gin.New()
	router.Use(gin.Recovery(), logger.GinMiddleware("Kong-Request-Id"))

	router.GET("_health", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	for i := range routes {
		var handlers []gin.HandlerFunc
		for _, h := range routes[i].middlewares {
			handlers = append(handlers, h)
		}
		handlers = append(handlers, routes[i].handler)
		router.Handle(routes[i].method, routes[i].endpoint, handlers...)
	}

	return router
}

func waitShutDown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	<-quit
	logger.Entry().Warning("Gracefully Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Entry().Fatal("Server Shutdown: ", err)
	}

	logger.Entry().Warning("Server exiting")
}
