package auth

import (
	"context"
	"log/slog"

	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"github.com/guluzadehh/bookapp/services/auth/internal/domain/models"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
)

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (*models.User, error)
}

type RoleProvider interface {
	GetRoleById(ctx context.Context, id int64) (*models.Role, error)
}

type AuthService struct {
	*services.Service
	config       *config.Config
	userProvider UserProvider
	roleProvider RoleProvider
}

func New(
	log *slog.Logger,
	config *config.Config,
	userProvider UserProvider,
	roleProvider RoleProvider,
) *AuthService {
	return &AuthService{
		Service:      services.New(log),
		config:       config,
		userProvider: userProvider,
		roleProvider: roleProvider,
	}
}
