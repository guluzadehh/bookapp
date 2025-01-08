package services

import (
	"errors"
	"log/slog"
)

var (
	ErrEmailExists  = errors.New("email is already taken")
	ErrRoleNotFound = errors.New("role doesn't exist")
)

type Service struct {
	Log *slog.Logger
}

func New(log *slog.Logger) *Service {
	return &Service{
		Log: log,
	}
}
