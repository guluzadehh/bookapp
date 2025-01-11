package auth

import (
	"context"
)

func (s *AuthService) Logout(ctx context.Context, accessStr, refreshStr string) {
	go s.blockToken(context.WithoutCancel(ctx), accessStr)
	go s.blockToken(context.WithoutCancel(ctx), refreshStr)
}
