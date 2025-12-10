package types

import (
	"pannypal/internal/common/enum"

	"github.com/google/uuid"
)

type UserWithAuth struct {
	ID       uuid.UUID     `json:"id" validate:"required"`
	Email    string        `json:"email" validate:"required,email"`
	UserType enum.UserType `json:"user_type" validate:"required,oneof=saas lite"`
	IsVerif  bool          `json:"is_verif" validate:"omitempty"`
	TeamId   string        `json:"team_id" validate:"omitempty"`
	RoleId   string        `json:"role_id" validate:"omitempty"`
}
