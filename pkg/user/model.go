package user

import "time"

type User struct {
	ID           string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Username     string    `json:"username,omitempty" bson:"username,omitempty"`
	Email        string    `json:"email,omitempty" bson:"email,omitempty"`
	Password     string    `json:"password,omitempty" bson:"password,omitempty"`
	CreatedAt    time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	Role         string    `json:"role,omitempty" bson:"role,omitempty"`
	Validity     string    `json:"validity,omitempty" bson:"validity,omitempty"`
	Code         string    `json:"code,omitempty" bson:"code,omitempty"`
	JoinedSalons []string  `json:"joinedSalons,omitempty" bson:"joinedRooms,omitempty"`
}

type UserValidation struct {
	Validation bool `json:"validation,omitempty"`
}

type UserUpdate struct {
	Username string `json:"username,omitempty"`
}

type PasswordUpdate struct {
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}
