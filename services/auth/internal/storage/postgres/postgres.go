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

func (s *Storage) CreateUser(ctx context.Context, email, password string, roleId int64) (*models.User, error) {
	const op = "storage.postgres.CreateUser"

	var user models.User

	const query = `
		INSERT INTO users(email, password, role_id)
		VALUES ($1, $2, $3)
		RETURNING id, email, password, role_id, created_at, updated_at, is_active;
	`

	var p string

	err := s.db.QueryRowContext(ctx, query, email, password, roleId).Scan(
		&user.Id,
		&user.Email,
		&p,
		&user.RoleId,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		if postgresErr, ok := err.(*pq.Error); ok {
			if postgresErr.Code.Name() == "unique_violation" {
				return nil, storage.UserExists
			}

			if postgresErr.Code.Name() == "foreign_key_violation" {
				return nil, storage.RoleNotFound
			}
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user.Password = []byte(p)

	return &user, nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "storage.postgres.UserByEmail"

	const query = `
		SELECT id, email, password, role_id, created_at, updated_at, is_active 
		FROM users
		WHERE email = $1 AND is_active = true;
	`

	var user models.User

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.RoleId,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.UserNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

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
