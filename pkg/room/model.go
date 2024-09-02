package room

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type RoomModel struct {
	ID          primitive.ObjectID ` bson:"_id,omitempty"`
	Name        string             ` bson:"name,omitempty"`
	Description string             ` bson:"description,omitempty"`
	Creator     string             ` bson:"creator,omitempty"`
	Members     []string           ` bson:"members,omitempty"`
	Hashtags    []string           ` bson:"hashtags,omitempty"`
	Messages    []string           ` bson:"messages,omitempty"`
	CreatedAt   time.Time          ` bson:"createdAt,omitempty"`
	UpdatedAt   time.Time          ` bson:"updatedAt,omitempty"`
}

func ModelToEntity(room *RoomModel) *RoomEntity {
	return &RoomEntity{
		ID:          room.ID.Hex(),
		Name:        room.Name,
		Description: room.Description,
		Creator:     room.Creator,
		Members:     room.Members,
		Hashtags:    room.Hashtags,
		Messages:    room.Messages,
		CreatedAt:   room.CreatedAt.String(),
		UpdatedAt:   room.UpdatedAt.String(),
	}
}

func EntityToModel(room *RoomEntity) *RoomModel {
	return &RoomModel{
		ID:          stringToObjectID(room.ID),
		Name:        room.Name,
		Description: room.Description,
		Creator:     room.Creator,
		Members:     room.Members,
		Hashtags:    room.Hashtags,
		Messages:    room.Messages,
		CreatedAt:   parseTime(room.CreatedAt),
		UpdatedAt:   parseTime(room.UpdatedAt),
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
