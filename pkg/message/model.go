package message

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RoomID    string             `json:"roomId,omitempty" bson:"roomId,omitempty"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty"`
	UserID    string             `json:"userId,omitempty" bson:"userId,omitempty"`
	Content   string             `json:"content,omitempty" bson:"content,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	// TODO: Add UpdateAt field
	//UpdateAt time.Time `json:"updateAt,omitempty" bson:"updateAt,omitempty"`
}
