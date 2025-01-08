package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"github.com/guluzadehh/bookapp/services/auth/internal/domain/models"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
	"github.com/guluzadehh/bookapp/services/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserSaver interface {
	CreateUser(ctx context.Context, email, password string) (*models.User, error)
}

type UserService struct {
	*services.Service
	config    *config.Config
	userSaver UserSaver
}

func New(
	log *slog.Logger,
	config *config.Config,
	userSaver UserSaver,
) *UserService {
	return &UserService{
		Service:   services.New(log),
		config:    config,
		userSaver: userSaver,
	}
}

func (s *UserService) CreateUser(
	ctx context.Context,
	email, password string,
) (*models.User, error) {
	const op = "services.user.CreateUser"

	log := sl.Init(s.Log, op, requestidmdw.GetId(ctx))

	var cost int = 14
	if s.config.Env == "local" {
		cost = 4
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Error("failed to hash password", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userSaver.CreateUser(ctx, email, string(bytes))
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
