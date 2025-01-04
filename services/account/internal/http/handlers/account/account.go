package accounthttp

import (
	"context"
	"log/slog"

	"github.com/guluzadehh/bookapp/services/account/internal/domain/models"
	httpbase "github.com/guluzadehh/bookapp/services/account/internal/http/base"
)

type AccountService interface {
	InitAccount(ctx context.Context, email, password string) (*models.User, error)
	SetLog(log *slog.Logger)
}

type AccountHandler struct {
	*httpbase.Handler
	srvc AccountService
}

func New(log *slog.Logger, srvc AccountService) *AccountHandler {
	return &AccountHandler{
		Handler: httpbase.New(log),
		srvc:    srvc,
	}
}
