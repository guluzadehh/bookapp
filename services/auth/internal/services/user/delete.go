package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/http/middlewares/authmdw"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
	"github.com/guluzadehh/bookapp/services/auth/internal/storage"
)

func (s *UserService) Delete(ctx context.Context, email string) error {
	const op = "services.user.Delete"

	log := sl.Init(s.Log, op, requestidmdw.GetId(ctx))

	user := authmdw.User(ctx)
	if user == nil || (user.Email != email && !user.Role.IsAdmin()) {
		log.Warn("unauthorized attempt to make user inactive")
		return fmt.Errorf("%s: %w", op, services.ErrUnauthorized)
	}

	if err := s.userStorage.DeleteUserByEmail(ctx, email); err != nil {
		if errors.Is(err, storage.UserNotFound) {
			log.Info("user not found")
			return fmt.Errorf("%s: %w", op, services.ErrUserNotFound)
		}

		log.Error("failed to make user inactive", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		if err := s.tokenBlacklist.BlacklistUser(context.WithoutCancel(ctx), email, s.config.JWT.Access.Expire); err != nil {
			log.Error("failed to blacklist user for token", sl.Err(err))
			return
		}

		log.Info("user has been blacklisted")
	}()

	return nil
}
