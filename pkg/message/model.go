package message

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// MessageEntity struct
type MessageModel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	RoomID    string             `bson:"roomId,omitempty"`
	Username  string             `bson:"username,omitempty"`
	UserID    string             `bson:"userId,omitempty"`
	Content   string             `bson:"content,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty"`
	// TODO: Add UpdateAt field
	//UpdateAt time.Time `json:"updateAt,omitempty" bson:"updateAt,omitempty"`
}

// MessageEntity struct
func ModelToEntity(message *MessageModel) *MessageEntity {
	return &MessageEntity{
		ID:        message.ID.Hex(),
		RoomID:    message.RoomID,
		Username:  message.Username,
		UserID:    message.UserID,
		Content:   message.Content,
		CreatedAt: message.CreatedAt.String(),
	}
}

// MessageEntity struct
func EntityToModel(message *MessageEntity) *MessageModel {
	return &MessageModel{
		ID:        stringToObjectID(message.ID),
		RoomID:    message.RoomID,
		Username:  message.Username,
		UserID:    message.UserID,
		Content:   message.Content,
		CreatedAt: parseTime(message.CreatedAt),
	}
}

// parseTime parses a time string and returns a time.Time object
func parseTime(timeStr string) time.Time {
	parsedTime, _ := time.Parse(time.RFC3339, timeStr)
	return parsedTime
}

// convert string to Object id
func stringToObjectID(id string) primitive.ObjectID {
	objectID, _ := primitive.ObjectIDFromHex(id)
	return objectID
}
