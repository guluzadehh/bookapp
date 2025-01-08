package storage

import "errors"

var (
	UserExists   = errors.New("user already exists")
	RoleNotFound = errors.New("role with that name doesn't exist")
)
