package main

import (
	"github.com/linuxfight/deepseek4free/internal/config"
	"github.com/linuxfight/deepseek4free/internal/stub"
	"github.com/linuxfight/deepseek4free/pkg/api"
	"github.com/linuxfight/deepseek4free/pkg/solver"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	wasmSolver, err := solver.New()
	if err != nil {
		panic(err)
	}

	apiClient := api.New(wasmSolver, cfg.RangersId, cfg.ApiKey)

	server := stub.New(cfg, apiClient, wasmSolver)

	/*
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {

		}()

		<-sigChan

		server.Stop()
	*/

	println("server started on :9090")
	if err := server.Listen(); err != nil {
		server.Stop()
		println(err.Error())
	}
}
