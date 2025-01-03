package account

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/account/internal/config"
	"github.com/guluzadehh/bookapp/services/account/internal/domain/models"
	"github.com/guluzadehh/bookapp/services/account/internal/services"
	"github.com/guluzadehh/bookapp/services/account/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type AccountSaver interface {
	CreateInitAccount(ctx context.Context, email, password string) (*models.User, error)
}

type AccountService struct {
	*services.Service
	config       *config.Config
	accountSaver AccountSaver
}

func New(
	log *slog.Logger,
	config *config.Config,
	accountSaver AccountSaver,
) *AccountService {
	return &AccountService{
		Service:      services.New(log),
		config:       config,
		accountSaver: accountSaver,
	}
}

func (s *AccountService) InitAccount(
	ctx context.Context,
	email, password string,
) (*models.User, error) {
	const op = "services.auth.CreateUser"

	log := s.Log.With(slog.String("op", op))

	var cost int = 14
	if s.config.Env == "local" {
		cost = 4
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Error("failed to hash password", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.accountSaver.CreateInitAccount(ctx, email, string(bytes))
	if err != nil {
		if errors.Is(err, storage.UserExists) {
			log.Info("email is taken")
			return nil, services.ErrEmailExists
		}

		log.Error("couldn't save the user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user has been created")

	return user, nil
}
