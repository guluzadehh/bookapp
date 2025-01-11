package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/lib/jwt"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
	"github.com/guluzadehh/bookapp/services/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func (s *AuthService) Authenticate(ctx context.Context, email string, password string) (string, string, error) {
	const op = "services.auth.Authenticate"

	log := s.Log.With(slog.String("op", op))

	user, err := s.userProvider.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.UserNotFound) {
			log.Warn("user doesn't exist")
			return "", "", fmt.Errorf("%s: %w", op, services.ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Warn("invalid credentials", slog.Int64("user_id", user.Id))
			return "", "", fmt.Errorf("%s: %w", op, services.ErrInvalidCredentials)
		}

		log.Error("failed to compare passwords", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	userRole, err := s.roleProvider.GetRoleById(ctx, user.RoleId)
	if err != nil {
		log.Error("failed to get user role", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	access, err := jwt.MakeAccess(user.Email, userRole.Name, s.config)
	if err != nil {
		log.Error("failed to generate access token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refresh, err := jwt.MakeRefresh(user.Email, userRole.Name, s.config)
	if err != nil {
		log.Error("failed to generate refresh token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("access and refresh tokens have been created")

	return access, refresh, nil
}
