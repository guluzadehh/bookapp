package account

import (
	"log/slog"

	"github.com/guluzadehh/bookapp/services/account/internal/config"
	"github.com/guluzadehh/bookapp/services/account/internal/services"
)

type AccountService struct {
	*services.Service
	config *config.Config
}

func New(log *slog.Logger, config *config.Config) *AccountService {
	return &AccountService{
		Service: services.New(log),
		config:  config,
	}
}
