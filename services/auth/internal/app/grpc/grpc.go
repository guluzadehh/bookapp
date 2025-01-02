package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"google.golang.org/grpc"
)

type GrpcApp struct {
	log        *slog.Logger
	grpcServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	config *config.Config,
) *GrpcApp {
	server := grpc.NewServer()

	return &GrpcApp{
		log:        log,
		grpcServer: server,
		port:       config.GRPCServer.Port,
	}
}

func (a *GrpcApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *GrpcApp) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("starting grpc server", slog.String("addr", l.Addr().String()))

	if err := a.grpcServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *GrpcApp) Stop() {
	const op = "grpcapp.Stop"

	log := a.log.With(slog.String("op", op))

	a.grpcServer.GracefulStop()
	log.Info("grpc server has been gracefully stopped")
}
