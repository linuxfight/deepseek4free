package main

import (
	"github.com/linuxfight/deepseek4free/internal/application"
	"github.com/linuxfight/deepseek4free/pkg/solver"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger.Info("logger initialized")

	wasmSolver, err := solver.New()
	if err != nil {
		panic(err)
	}

	logger.Info("wasm solver initialized")

	app := application.New(wasmSolver, logger)
	app.Init()

	logger.Info("app initialized")

	// TODO: context, get all messages from dialog and give llm only the final messages, make it answer the last question
	// TODO 2: add logs, traces, metrics

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("application started on port 8080")
		err := app.Start()
		if err != nil {
			stopErr := app.Stop()
			if stopErr != nil {
				return
			}
			panic(err)
		}
	}()

	<-sigChan

	err = app.Stop()
	if err != nil {
		panic(err)
	}
}
