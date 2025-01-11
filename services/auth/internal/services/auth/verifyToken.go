package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/guluzadehh/bookapp/services/auth/internal/lib/jwt"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
)

func (s *AuthService) VerifyToken(
	ctx context.Context,
	tokenStr string,
) (*jwt.AuthClaims, error) {
	const op = "services.auth.VerifyToken"

	// jwt verification of token
	token, err := jwt.Verify(tokenStr, s.config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, services.ErrInvalidToken)
	}

	if blacklisted, err := s.tokenBlacklist.TokenInBlacklist(ctx, tokenStr); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	} else if blacklisted {
		return nil, fmt.Errorf("%s: %w", op, services.ErrInvalidToken)
	}

	claims, ok := token.Claims.(*jwt.AuthClaims)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, errors.New("failed to get jwt claims"))
	}

	if blacklisted, err := s.tokenBlacklist.UserInBlacklist(ctx, claims.Email); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	} else if blacklisted {
		return nil, fmt.Errorf("%s: %w", op, services.ErrInvalidToken)
	}

	if blacklisted, err := s.tokenBlacklist.RoleInBlacklist(ctx, claims.Role); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	} else if blacklisted {
		return nil, fmt.Errorf("%s: %w", op, services.ErrInvalidToken)
	}

	return claims, nil
}
