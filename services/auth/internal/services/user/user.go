package user

import (
	"log/slog"

	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
)

type UserService struct {
	*services.Service
	config *config.Config
}

func New(log *slog.Logger, config *config.Config) *UserService {
	return &UserService{
		Service: services.New(log),
		config:  config,
	}
}
