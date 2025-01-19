package app

import (
	"log/slog"

	"github.com/guluzadehh/bookapp/pkg/sl"
	grpcapp "github.com/guluzadehh/bookapp/services/auth/internal/app/grpc"
	httpapp "github.com/guluzadehh/bookapp/services/auth/internal/app/http"
	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"github.com/guluzadehh/bookapp/services/auth/internal/services/auth"
	"github.com/guluzadehh/bookapp/services/auth/internal/services/user"
	"github.com/guluzadehh/bookapp/services/auth/internal/storage/postgres"
	"github.com/guluzadehh/bookapp/services/auth/internal/storage/redis"
)

type App struct {
	log      *slog.Logger
	httpApp  *httpapp.HttpApp
	grpcApp  *grpcapp.GrpcApp
	postgres *postgres.Storage
}

func New(log *slog.Logger, config *config.Config) *App {
	log.Info("connecting to postgres")
	pgStorage, err := postgres.New(config.Postgres.DSN(config.Postgres.Options))
	if err != nil {
		panic(err)
	}

	log.Info("connecting to redis")
	redisStorage, err := redis.New(config)
	if err != nil {
		panic(err)
	}

	userService := user.New(log, config, pgStorage, pgStorage, redisStorage)
	authService := auth.New(log, config, pgStorage, pgStorage, redisStorage)

	httpApp := httpapp.New(log, config, userService, authService, authService, userService)
	grpcApp := grpcapp.New(log, config)

	return &App{
		log:      log,
		httpApp:  httpApp,
		grpcApp:  grpcApp,
		postgres: pgStorage,
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

	a.log.Info("closing postgres connection")
	if err := a.postgres.Close(); err != nil {
		a.log.Error("failed to close the postgres connection", sl.Err(err))
	}
	a.log.Info("postgres connection has been stopped")
}
