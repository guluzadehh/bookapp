package userhttp

import (
	"log/slog"

	httpbase "github.com/guluzadehh/bookapp/services/auth/internal/http/base"
)

type UserService interface{}

type UserHandler struct {
	*httpbase.Handler
	srvc UserService
}

func New(log *slog.Logger, srvc UserService) *UserHandler {
	return &UserHandler{
		Handler: httpbase.New(log),
		srvc:    srvc,
	}
}
