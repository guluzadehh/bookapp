package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/guluzadehh/bookapp/services/auth/internal/app"
	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"github.com/joho/godotenv"
)

const (
	env_local = "local"
	env_dev   = "dev"
	env_prod  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed to load .env file: %s\n", err.Error())
	}

	config := config.MustLoad()

	log := setupLogger(config.Env).With(slog.String("app", "auth"))

	log.Info("building app", slog.String("env", config.Env))
	app := app.New(log, config)

	log.Info("starting app")
	app.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("gracefully stopping app")
	app.Stop()

	log.Info("app has been gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case env_local, env_dev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case env_prod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
