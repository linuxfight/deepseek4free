package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/linuxfight/deepseek4free/internal/application"
	"github.com/linuxfight/deepseek4free/internal/kv"
	"github.com/linuxfight/deepseek4free/pkg/solver"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger.Info("logger initialized")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	cache, err := kv.New(redisAddr)
	if err != nil {
		logger.Fatal("failed to initialize cache", zap.Error(err))
	}

	wasmSolver, err := solver.New()
	if err != nil {
		logger.Fatal("failed to initialize wasm solver", zap.Error(err))
	}

	logger.Info("wasm solver initialized")

	app := application.New(wasmSolver, logger, cache)
	app.Init()

	logger.Info("app initialized")

	// TODO: add logs, traces, metrics

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("application started on port 8080")
		err := app.Start()
		if err != nil {
			stopErr := app.Stop()
			if stopErr != nil {
				logger.Fatal("failed to stop application", zap.Error(stopErr))
			}
			logger.Fatal("application stopped", zap.Error(err))
		}
	}()

	<-sigChan

	err = app.Stop()
	if err != nil {
		logger.Fatal("failed to stop application", zap.Error(err))
	}

	logger.Info("application stopped")
}
