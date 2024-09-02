package message

import (
	"chat-app/pkg/room"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// MessageRepository defines the methods that a message repository should implement
type MessageRepository interface {
	CreateMessage(ctx context.Context, message *MessageEntity) (*MessageEntity, error)
	GetMessages(ctx context.Context, roomID string) ([]*MessageEntity, error)
	GetMessage(ctx context.Context, messageID string) (*MessageEntity, error)
	DeleteMessage(ctx context.Context, messageID string) error
}

// messageRepository is a struct that embeds the collection of messages and rooms
type messageRepository struct {
	collectionMessage *mongo.Collection
	collectionRoom    *mongo.Collection
}

// NewMessageRepository creates a new instance of MessageRepository
func NewMessageRepository(collectionMessage *mongo.Collection, collectionRoom *mongo.Collection) MessageRepository {
	return &messageRepository{collectionMessage: collectionMessage, collectionRoom: collectionRoom}
}

// CreateMessage creates a new message in the database
func (r *messageRepository) CreateMessage(ctx context.Context, message *MessageEntity) (*MessageEntity, error) {
	// add others required fields by default to the message
	messageModel := &MessageModel{
		ID:        primitive.NewObjectID(),
		RoomID:    message.RoomID,
		Username:  message.Username,
		UserID:    message.UserID,
		Content:   message.Content,
		CreatedAt: time.Now(),
	}
	// insert message into the message collection
	_, err := r.collectionMessage.InsertOne(ctx, messageModel)
	if err != nil {
		return nil, err
	}

	// update the room collection with the new message
	// convert roomID to string
	roomPrimitiveID, err := primitive.ObjectIDFromHex(messageModel.RoomID)
	if err != nil {
		return nil, err
	}
	_, err = r.collectionRoom.UpdateOne(ctx, bson.D{{"_id", roomPrimitiveID}}, bson.D{{"$push", bson.D{{"messages", messageModel.ID.Hex()}}}})
	if err != nil {
		// delete the message if the update fails
		_, err = r.collectionMessage.DeleteOne(ctx, bson.D{{"_id", messageModel.ID}})
		if err != nil {
			return nil, err
		}
		return nil, err

	}
	return ModelToEntity(messageModel), nil
}

// GetMessages retrieves all messages from a room
func (r *messageRepository) GetMessages(ctx context.Context, roomID string) ([]*MessageEntity, error) {
	// convert roomID to ObjectID
	roomIDObjectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, err
	}
	// retrieve the room from the database
	var roomRetrieved room.RoomModel
	err = r.collectionRoom.FindOne(ctx, bson.D{{"_id", roomIDObjectID}}).Decode(&roomRetrieved)
	if err != nil {
		fmt.Println("Error 1 : ", err)
		return nil, err

	}

	var messages []*MessageModel
	for _, messageId := range roomRetrieved.Messages {
		var messageRetrieved *MessageModel
		// change ID string into ObjectID
		messageObjectID, err := primitive.ObjectIDFromHex(messageId)
		if err != nil {
			// debug
			return nil, err
		}

		err = r.collectionMessage.FindOne(ctx, bson.D{{"_id", messageObjectID}}).Decode(&messageRetrieved)
		if err != nil {
			return nil, err
		}
		messages = append(messages, messageRetrieved)
	}

	// convert the messages to entities
	var messagesEntities []*MessageEntity
	for _, message := range messages {
		messagesEntities = append(messagesEntities, ModelToEntity(message))
	}
	return messagesEntities, nil
}

// GetMessage converts a string to an ObjectID
func (r *messageRepository) GetMessage(ctx context.Context, messageID string) (*MessageEntity, error) {
	var message MessageModel
	// convert messageID to ObjectID
	messageIDObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, err
	}
	err = r.collectionMessage.FindOne(ctx, bson.D{{"_id", messageIDObjectID}}).Decode(&message)
	if err != nil {
		return nil, err
	}
	return ModelToEntity(&message), nil

}

// DeleteMessage deletes a message from the database
func (r *messageRepository) DeleteMessage(ctx context.Context, messageID string) error {
	// convert messageID to ObjectID
	messageIDObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return err
	}

	// Define the filter to find the message with the given ID
	filter := bson.M{"_id": messageIDObjectID}

	// Delete the message from the messages collection
	_, err = r.collectionMessage.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// Update the rooms collection to remove the message ID from the list of messages
	update := bson.M{"$pull": bson.M{"messages": messageID}}

	// Update all rooms to remove the message ID from the list of messages
	_, err = r.collectionRoom.UpdateMany(ctx, bson.M{}, update)
	if err != nil {
		// If the update fails, reinsert the message ID into the list of messages
		_, err := r.collectionRoom.UpdateMany(ctx, bson.M{}, bson.M{"$push": bson.M{"messages": messageID}})
		if err != nil {
			return err
		}
		return err
	}

	return nil
}
