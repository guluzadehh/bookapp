package models

import "time"

type Role struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Role) IsAdmin() bool {
	return r.Name == "admin"
}

func (r *Role) IsUser() bool {
	return r.Name == "user"
}
