package application

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/linuxfight/deepseek4free/internal/kv"
	"github.com/linuxfight/deepseek4free/internal/loggerware"
	"github.com/linuxfight/deepseek4free/internal/serializer"
	"github.com/linuxfight/deepseek4free/pkg/solver"
	"go.uber.org/zap"
)

type Instance struct {
	router *echo.Echo
	solver *solver.Instance
	cache  *kv.Instance
}

func (i *Instance) Init() {
	i.router.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "started")
	})
	i.router.GET("/models", i.models)
	i.router.POST("/chat/completions", i.chat)
}

func (i *Instance) Start() error {
	return i.router.Start(":8080")
}

func (i *Instance) Stop() error {
	i.cache.Close()
	i.solver.Close()
	return i.router.Shutdown(context.Background())
}

func New(solver *solver.Instance, logger *zap.Logger, cache *kv.Instance) *Instance {
	router := echo.New()

	router.HideBanner = true
	router.HidePort = true
	router.Use(middleware.Recover())
	router.Use(loggerware.ZapLogger(logger))
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	router.JSONSerializer = serializer.New()

	return &Instance{
		router: router,
		solver: solver,
		cache:  cache,
	}
}
