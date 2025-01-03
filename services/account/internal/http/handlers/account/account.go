package accounthttp

import (
	"log/slog"

	httpbase "github.com/guluzadehh/bookapp/services/account/internal/http/base"
)

type AccountService interface {
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
