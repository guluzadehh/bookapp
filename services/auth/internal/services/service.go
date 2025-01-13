package services

import (
	"errors"
	"log/slog"
)

var (
	ErrEmailExists        = errors.New("email is already taken")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrRoleNotFound       = errors.New("role doesn't exist")
	ErrInvalidToken       = errors.New("invalid token")
)

type Service struct {
	Log *slog.Logger
}

func New(log *slog.Logger) *Service {
	return &Service{
		Log: log,
	}
}
