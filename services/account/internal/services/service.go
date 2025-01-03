package services

import (
	"errors"
	"log/slog"
)

var (
	ErrEmailExists = errors.New("email is already taken")
)

type Service struct {
	Log *slog.Logger
}

func New(log *slog.Logger) *Service {
	return &Service{
		Log: log,
	}
}

func (s *Service) SetLog(log *slog.Logger) {
	s.Log = log
}
