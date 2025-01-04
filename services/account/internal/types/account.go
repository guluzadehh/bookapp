package types

import (
	"time"

	"github.com/guluzadehh/bookapp/services/account/internal/domain/models"
)

type AccountView struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

func NewAccount(u *models.User) *AccountView {
	if u == nil {
		return nil
	}

	return &AccountView{
		Id:        u.Id,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsActive:  u.IsActive,
	}
}
