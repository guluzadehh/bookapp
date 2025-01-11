package types

import (
	"time"

	"github.com/guluzadehh/bookapp/services/auth/internal/domain/models"
)

type UserView struct {
	Id        int64         `json:"id"`
	Email     string        `json:"email"`
	Role      *userRoleView `json:"role,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	IsActive  bool          `json:"is_active"`
}

func NewUser(u *models.User) *UserView {
	if u == nil {
		return nil
	}

	return &UserView{
		Id:        u.Id,
		Email:     u.Email,
		Role:      newUserRoleView(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsActive:  u.IsActive,
	}
}

type userRoleView struct {
	Name string `json:"name"`
}

func newUserRoleView(r *models.Role) *userRoleView {
	if r == nil {
		return nil
	}

	return &userRoleView{
		Name: r.Name,
	}
}
