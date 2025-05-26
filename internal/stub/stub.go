package stub

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/linuxfight/deepseek4free/internal/config"
	"github.com/linuxfight/deepseek4free/internal/stub/gen"
	"github.com/linuxfight/deepseek4free/pkg/api"
	"github.com/linuxfight/deepseek4free/pkg/solver"
	"google.golang.org/grpc"
	"net"
)

type Stub struct {
	gen.UnimplementedDeepseekApiServer

	config     *config.Config
	server     *grpc.Server
	api        *api.Client
	wasmSolver *solver.Solver
}

func (stub *Stub) Listen() error {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		return err
	}
	if err := stub.server.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (stub *Stub) Stop() {
	stub.server.GracefulStop()
	stub.wasmSolver.Close()
}

func New(config *config.Config, api *api.Client, wasmSolver *solver.Solver) *Stub {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(recovery.UnaryServerInterceptor()),
	)
	stub := &Stub{
		api:        api,
		wasmSolver: wasmSolver,
		config:     config,
		server:     server,
	}
	gen.RegisterDeepseekApiServer(server, stub)
	return stub
}
