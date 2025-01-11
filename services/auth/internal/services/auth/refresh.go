package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/lib/jwt"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
)

func (s *AuthService) RefreshToken(ctx context.Context, refreshStr, oldAccessStr string) (string, error) {
	const op = "services.auth.RefreshToken"

	log := sl.Init(s.Log, op, requestidmdw.GetId(ctx))

	claims, err := s.VerifyToken(ctx, refreshStr)
	if err != nil {
		if errors.Is(err, services.ErrInvalidToken) {
			log.Info("refresh token is invalid", slog.String("invalid_refresh_token", refreshStr), sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, err)
		}

		log.Error("error while getting the subject from refresh token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	access, err := jwt.MakeAccess(claims.Email, claims.Role, s.config)
	if err != nil {
		log.Error("can't create jwt access token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	log.Info("access token have been created")

	go func() {
		if oldAccessStr == "" {
			return
		}

		claims, err := s.VerifyToken(context.Background(), oldAccessStr)
		if err != nil {
			if errors.Is(err, services.ErrInvalidToken) {
				return
			}

			// Unknown token verification error
			log.Error("token verification error", sl.Err(err))
			return
		}

		expiration, err := claims.GetExpirationTime()
		if err != nil {
			log.Error("failed to get token exp time", sl.Err(err))
			return
		}

		if err := s.tokenBlacklist.BlacklistToken(context.Background(), oldAccessStr, time.Until(expiration.Time)); err != nil {
			log.Error("failed to add old token to blacklist", sl.Err(err))
			// TODO
		}

		log.Info("token has been blacklisted")
	}()

	return access, nil
}
