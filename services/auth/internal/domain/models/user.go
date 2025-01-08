package models

import "time"

type User struct {
	Id        int64
	Email     string
	Password  []byte
	RoleId    int64
	Role      *Role
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}
