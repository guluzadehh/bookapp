package storage

import "errors"

var (
	UserExists   = errors.New("user already exists")
	UserNotFound = errors.New("user not found")
	RoleNotFound = errors.New("role with that name doesn't exist")
)
