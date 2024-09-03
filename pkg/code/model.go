package code

import "time"

// CodeModel to be used only once by user in registration
type CodeModel struct {
	Code      string    `bson:"code,omitempty"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
	IsUsed    bool      `bson:"isUsed"`
}
