package application

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/linuxfight/deepseek4free/internal/loggerware"
	"github.com/linuxfight/deepseek4free/internal/serializer"
	"github.com/linuxfight/deepseek4free/pkg/solver"
	"go.uber.org/zap"
	"net/http"
)

type Instance struct {
	router *echo.Echo

	solver *solver.Solver
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
	return i.router.Shutdown(context.Background())
}

func New(solver *solver.Solver, logger *zap.Logger) *Instance {
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
	}
}
