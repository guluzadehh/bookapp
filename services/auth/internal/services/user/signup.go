package user

import (
	"context"
	"log/slog"

	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/domain/models"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
)

func (s *UserService) Signup(
	ctx context.Context,
	email, password string,
) (*models.User, error) {
	const op = "services.user.Signup"

	log := s.Log.With(slog.String("op", op))

	userRole, err := s.roleProvider.GetRoleByName(ctx, "user")
	if err != nil {
		log.Error("failed to get the user role", sl.Err(err))
		return nil, services.ErrRoleNotFound
	}

	user, err := s.createUser(ctx, email, password, userRole.Id)
	if err != nil {
		return nil, err
	}

	user.Role = userRole

	return user, nil
}
