package code

import "time"

// code to be used only once by user in registration

type Code struct {
	Code      string    `json:"code,omitempty" bson:"code,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	IsUsed    bool      `json:"isUsed,omitempty" bson:"isUsed"`
}
