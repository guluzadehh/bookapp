package authhttp

import (
	"context"
	"log/slog"

	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	httpbase "github.com/guluzadehh/bookapp/services/auth/internal/http/base"
	"github.com/guluzadehh/bookapp/services/auth/internal/lib/jwt"
)

type AuthService interface {
	Authenticate(ctx context.Context, email, password string) (access string, refresh string, err error)
	VerifyToken(ctx context.Context, token string) (*jwt.AuthClaims, error)
	RefreshToken(ctx context.Context, refresh string, oldAccess string) (access string, err error)
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
