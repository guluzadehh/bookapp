package httpapp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guluzadehh/bookapp/pkg/http/middlewares/loggingmdw"
	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	authhttp "github.com/guluzadehh/bookapp/services/auth/internal/http/handlers/auth"
	userhttp "github.com/guluzadehh/bookapp/services/auth/internal/http/handlers/user"
	"github.com/guluzadehh/bookapp/services/auth/internal/http/middlewares/authmdw"
)

type HttpApp struct {
	log        *slog.Logger
	httpServer *http.Server
}

func New(
	log *slog.Logger,
	config *config.Config,
	userService userhttp.UserService,
	authService authhttp.AuthService,
	authMdwAuthService authhttp.AuthService,
	authMdwUserService authmdw.UserService,
) *HttpApp {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.HTTPServer.Port),
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}

	userHandler := userhttp.New(log, userService)
	authHandler := authhttp.New(log, config, authService)

	router := mux.NewRouter()
	router.Use(loggingmdw.Middleware(log))
	router.Use(requestidmdw.Middleware)

	api := router.PathPrefix("/api").Subrouter()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", userHandler.Signup).Methods("POST")
	auth.HandleFunc("/login", authHandler.Authenticate).Methods("POST")
	auth.HandleFunc("/refresh", authHandler.Refresh).Methods("POST")

	protectedAuth := auth.NewRoute().Subrouter()
	protectedAuth.Use(authmdw.Authorize(log, config, authMdwAuthService, authMdwUserService))


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
