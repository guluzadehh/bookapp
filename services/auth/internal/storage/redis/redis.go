package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/guluzadehh/bookapp/services/auth/internal/config"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	cli *redis.Client
}

func New(config *config.Config) (*Storage, error) {
	cli := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DefaultDB,
	})

	if err := cli.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return &Storage{
		cli: cli,
	}, nil
}

func (s *Storage) BlacklistToken(ctx context.Context, token string, expiry time.Duration) error {
	const op = "storage.redis.BlacklistToken"

	err := s.cli.Set(ctx, fmt.Sprintf("unactive:tokens:%s", token), true, expiry).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) BlacklistUser(ctx context.Context, email string, expiry time.Duration) error {
	const op = "storage.redis.BlacklistToken"

	err := s.cli.Set(ctx, fmt.Sprintf("unactive:users:%s", email), true, expiry).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) BlacklistRole(ctx context.Context, role string, expiry time.Duration) error {
	const op = "storage.redis.BlacklistToken"

	err := s.cli.Set(ctx, fmt.Sprintf("unactive:roles:%s", role), true, expiry).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) TokenInBlacklist(ctx context.Context, token string) (bool, error) {
	const op = "storage.redis.TokenInBlacklist"

	exists, err := s.cli.Exists(ctx, fmt.Sprintf("unactive:tokens:%s", token)).Result()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists > 0, nil
}

func (s *Storage) UserInBlacklist(ctx context.Context, email string) (bool, error) {
	const op = "storage.redis.UserInBlacklist"

	exists, err := s.cli.Exists(ctx, fmt.Sprintf("unactive:users:%s", email)).Result()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists > 0, nil
}

func (s *Storage) RoleInBlacklist(ctx context.Context, role string) (bool, error) {
	const op = "storage.redis.RoleInBlacklist"

	exists, err := s.cli.Exists(ctx, fmt.Sprintf("unactive:roles:%s", role)).Result()
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists > 0, nil
}
