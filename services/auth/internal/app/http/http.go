package httpapp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	loggingmdw "github.com/guluzadehh/bookapp/pkg/http/middlewares"
	"github.com/guluzadehh/bookapp/services/auth/internal/config"
)

type HttpApp struct {
	log        *slog.Logger
	httpServer *http.Server
}

func New(
	log *slog.Logger,
	config *config.Config,
) *HttpApp {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.HTTPServer.Port),
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}

	router := mux.NewRouter()
	router.Use(loggingmdw.LogRequests(log))

	server.Handler = router

	return &HttpApp{
		log:        log,
		httpServer: server,
	}
}

func (a *HttpApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *HttpApp) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(slog.String("op", op))

	log.Info("starting http server", slog.String("addr", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (a *HttpApp) Stop() error {
	const op = "httpapp.Stop"

	log := a.log.With("op", op)

	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("http server has been gracefully stopped")
	return nil
}