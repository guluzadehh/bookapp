package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/guluzadehh/bookapp/services/auth/internal/domain/models"
	"github.com/guluzadehh/bookapp/services/auth/internal/storage"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateUser(ctx context.Context, email, password string) (*models.User, error) {
	const op = "storage.postgres.CreateUser"

	var user models.User

	const query = `
		INSERT INTO users(email, password)
		VALUES ($1, $2)
		RETURNING id, email, password, created_at, updated_at, is_active;
	`

	var p string

	err := s.db.QueryRowContext(ctx, query, email, password).Scan(
		&user.Id,
		&user.Email,
		&p,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok && postgresErr.Code.Name() == "unique_violation" {
			return nil, storage.UserExists
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user.Password = []byte(p)

	return &user, nil
}

func (s *Storage) GetRoleById(ctx context.Context, roleId int64) (*models.Role, error) {
	const op = "storage.postgres.GetRoleByName"

	const query = `
		SELECT id, name, created_at, updated_at
		FROM roles
		WHERE id = $1; 
	`

	var role models.Role

	err := s.db.QueryRowContext(ctx, query, roleId).Scan(
		&role.Id,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.RoleNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}

func (s *Storage) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	const op = "storage.postgres.GetRoleByName"

	const query = `
		SELECT id, name, created_at, updated_at
		FROM roles
		WHERE name = $1; 
	`

	var role models.Role

	err := s.db.QueryRowContext(ctx, query, name).Scan(
		&role.Id,
		&role.Name,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.RoleNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &role, nil
}
