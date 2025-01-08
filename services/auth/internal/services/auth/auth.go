package auth

import (
	"log/slog"

	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
)

type AuthService struct {
	*services.Service
	config *config.Config
}

func New(
	log *slog.Logger,
	config *config.Config,
) *AuthService {
	return &AuthService{
		Service: services.New(log),
		config:  config,
	}
}
