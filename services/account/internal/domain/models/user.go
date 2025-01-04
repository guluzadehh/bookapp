package models

import "time"

type User struct {
	Id        int64
	Email     string
	Password  []byte
	FirstName string
	LastName  string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsActive  bool
}
