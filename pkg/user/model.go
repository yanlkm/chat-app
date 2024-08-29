package user

import "time"

type UserModel struct {
	ID           string    `bson:"_id,omitempty"`
	Username     string    `bson:"username,omitempty"`
	Email        string    `bson:"email,omitempty"`
	Password     string    `bson:"password,omitempty"`
	CreatedAt    time.Time `bson:"createdAt,omitempty"`
	UpdatedAt    time.Time `bson:"updatedAt,omitempty"`
	Role         string    `bson:"role,omitempty"`
	Validity     string    `bson:"validity,omitempty"`
	JoinedSalons []string  `bson:"joinedRooms,omitempty"`
}
