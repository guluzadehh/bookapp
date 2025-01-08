package userhttp

import (
	"context"
	"log/slog"

	"github.com/guluzadehh/bookapp/services/auth/internal/domain/models"
	httpbase "github.com/guluzadehh/bookapp/services/auth/internal/http/base"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password string) (*models.User, error)
}

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
