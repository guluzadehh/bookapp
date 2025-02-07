package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/guluzadehh/bookapp/pkg/http/middlewares/requestidmdw"
	"github.com/guluzadehh/bookapp/pkg/sl"
	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"github.com/guluzadehh/bookapp/services/auth/internal/domain/models"
	"github.com/guluzadehh/bookapp/services/auth/internal/services"
	"github.com/guluzadehh/bookapp/services/auth/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	UserByEmailWithRole(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, email, password string, roleId int64) (*models.User, error)
	DeleteUserByEmail(ctx context.Context, email string) error
}

type RoleProvider interface {
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
}

type TokenBlacklist interface {
	BlacklistUser(ctx context.Context, email string, expiry time.Duration) error
}

type UserService struct {
	*services.Service
	config         *config.Config
	userStorage    UserStorage
	roleProvider   RoleProvider
	tokenBlacklist TokenBlacklist
}

func New(
	log *slog.Logger,
	config *config.Config,
	userStorage UserStorage,
	roleProvider RoleProvider,
	tokenBlacklist TokenBlacklist,
) *UserService {
	return &UserService{
		Service:        services.New(log),
		config:         config,
		userStorage:    userStorage,
		roleProvider:   roleProvider,
		tokenBlacklist: tokenBlacklist,
	}
}

func (s *UserService) createUser(
	ctx context.Context,
	email, password string,
	roleId int64,
) (*models.User, error) {
	const op = "services.user.createUser"

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

	user, err := s.userStorage.CreateUser(ctx, email, string(bytes), roleId)
	if err != nil {
		if errors.Is(err, storage.UserExists) {
			log.Info("email is taken")
			return nil, services.ErrEmailExists
		}

		if errors.Is(err, storage.RoleNotFound) {
			log.Error("role doesn't exist", sl.Err(err))
			return nil, services.ErrRoleNotFound
		}

		log.Error("couldn't save the user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user has been created")

	return user, nil
}
