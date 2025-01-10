package redis

import (
	"context"

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
