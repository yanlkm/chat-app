package room

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Room struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Creator     string             `json:"creator,omitempty" bson:"creator,omitempty"`
	Members     []string           `json:"members,omitempty" bson:"members,omitempty"`
	Hashtags    []string           `json:"hashtags,omitempty" bson:"hashtags,omitempty"`
	Messages    []string           `json:"messages,omitempty" bson:"messages,omitempty"`
	CreatedAt   time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}
