package user

import "time"

type UserEntity struct {
	ID           string    `json:"_id,omitempty"`
	Username     string    `json:"username,omitempty"`
	Email        string    `json:"email,omitempty"`
	Password     string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
	Role         string    `json:"role,omitempty"`
	Validity     string    `json:"validity,omitempty"`
	Code         string    `json:"code,omitempty"`
	JoinedSalons []string  `json:"joinedSalons,omitempty"`
}

type UserValidationEntity struct {
	Validation bool `json:"validation,omitempty"`
}

type UserUpdateEntity struct {
	Username string `json:"username,omitempty"`
}

type PasswordUpdateEntity struct {
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}
