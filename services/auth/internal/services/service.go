package services

import "log/slog"

type Service struct {
	Log *slog.Logger
}

func New(log *slog.Logger) *Service {
	return &Service{
		Log: log,
	}
}
