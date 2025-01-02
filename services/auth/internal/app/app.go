package app

import (
	"log/slog"

	"github.com/guluzadehh/bookapp/pkg/sl"
	grpcapp "github.com/guluzadehh/bookapp/services/auth/internal/app/grpc"
	httpapp "github.com/guluzadehh/bookapp/services/auth/internal/app/http"
	"github.com/guluzadehh/bookapp/services/auth/internal/config"
)

type App struct {
	log     *slog.Logger
	httpApp *httpapp.HttpApp
	grpcApp *grpcapp.GrpcApp
}

func New(log *slog.Logger, config *config.Config) *App {
	httpApp := httpapp.New(log, config)
	grpcApp := grpcapp.New(log, config)

	return &App{
		log:     log,
		httpApp: httpApp,
		grpcApp: grpcApp,
	}
}

func (a *App) Start() {
	a.log.Info("starting http app")
	go a.httpApp.MustRun()

	a.log.Info("starting grpc app")
	go a.grpcApp.MustRun()
}

func (a *App) Stop() {
	a.log.Info("stopping http app")
	if err := a.httpApp.Stop(); err != nil {
		a.log.Error("failed t gracefully stop http server", sl.Err(err))
	}

	a.log.Info("stopping grpc app")
	a.grpcApp.Stop()
}
