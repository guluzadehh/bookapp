package auth

import (
	"context"
	"errors"
	"time"

	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
)

func (s *AuthService) blockToken(ctx context.Context, token string) {
	const op = "services.auth.BlockToken"

	log := sl.Init(s.Log, op, requestidmdw.GetId(ctx))

	if token == "" {
		return
	}

	claims, err := s.VerifyToken(ctx, token)
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

	if err := s.tokenBlacklist.BlacklistToken(ctx, token, time.Until(expiration.Time)); err != nil {
		log.Error("failed to add old token to blacklist", sl.Err(err))
		// TODO
		return
	}

	log.Info("token has been blacklisted")
}
