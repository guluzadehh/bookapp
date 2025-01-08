package authhttp

import (
	"context"
	"log/slog"

	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	httpbase "github.com/guluzadehh/bookapp/services/auth/internal/http/base"
)

type AuthService interface {
	Authenticate(ctx context.Context, email, password string) (access string, refresh string, err error)
}

type AuthHandler struct {
	*httpbase.Handler
	config *config.Config
	srvc   AuthService
}

func New(
	log *slog.Logger,
	config *config.Config,
	srvc AuthService,
) *AuthHandler {
	return &AuthHandler{
		Handler: httpbase.New(log),
		config:  config,
		srvc:    srvc,
	}
}
