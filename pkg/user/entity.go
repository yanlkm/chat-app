package user

import "time"

// UserEntity represents the structure of a user entity.
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

// UserValidationEntity  represents the login credentials provided by the user.
type UserValidationEntity struct {
	Validation bool `json:"validation,omitempty"`
}

// UserUpdateEntity represents the username update credentials provided by the user.
type UserUpdateEntity struct {
	Username string `json:"username,omitempty"`
}

// PasswordUpdateEntity represents the password update credentials provided by the user.
type PasswordUpdateEntity struct {
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}
