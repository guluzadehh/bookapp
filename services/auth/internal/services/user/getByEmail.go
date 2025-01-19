package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/domain/models"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
	"github.com/guluzadehh/bookapp/services/auth/internal/storage"
)

func (s *UserService) GetUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {
	const op = "services.user.GetUserByEmail"

	log := sl.Init(s.Log, op, requestidmdw.GetId(ctx))

	user, err := s.userStorage.UserByEmailWithRole(ctx, email)
	if err != nil {
		if errors.Is(err, storage.UserNotFound) {
			log.Info("user doesn't exist")
			return nil, fmt.Errorf("%s: %w", op, services.ErrUserNotFound)
		}

		log.Error("failed to get user from storage", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
