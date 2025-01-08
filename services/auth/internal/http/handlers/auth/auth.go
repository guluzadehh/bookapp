package authhttp

import (
	"log/slog"

	httpbase "github.com/guluzadehh/bookapp/services/auth/internal/http/base"
)

type AuthService interface{}

type AuthHandler struct {
	*httpbase.Handler
	srvc AuthService
}

func New(log *slog.Logger, srvc AuthService) *AuthHandler {
	return &AuthHandler{
		Handler: httpbase.New(log),
		srvc:    srvc,
	}
}
