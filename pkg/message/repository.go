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
	CreateMessage(ctx context.Context, message *Message) (*Message, error)
	GetMessages(ctx context.Context, roomID primitive.ObjectID) ([]*Message, error)
	GetMessage(ctx context.Context, messageID primitive.ObjectID) (*Message, error)
	DeleteMessage(ctx context.Context, messageID primitive.ObjectID) error
}

type messageRepository struct {
	collectionMessage *mongo.Collection
	collectionRoom    *mongo.Collection
}

// NewMessageRepository creates a new instance of MessageRepository
func NewMessageRepository(collectionMessage *mongo.Collection, collectionRoom *mongo.Collection) MessageRepository {
	return &messageRepository{collectionMessage: collectionMessage, collectionRoom: collectionRoom}
}

// CreateMessage creates a new message in the database
func (r *messageRepository) CreateMessage(ctx context.Context, message *Message) (*Message, error) {
	// add others required fields by default
	message = &Message{
		ID:        primitive.NewObjectID(),
		RoomID:    message.RoomID,
		Username:  message.Username,
		Content:   message.Content,
		CreatedAt: time.Now(),
	}
	// insert message into the message collection
	_, err := r.collectionMessage.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	// update the room collection with the new message
	// convert roomID to primitive.ObjectID
	roomPrimitiveID, err := primitive.ObjectIDFromHex(message.RoomID)
	if err != nil {
		return nil, err
	}
	_, err = r.collectionRoom.UpdateOne(ctx, bson.D{{"_id", roomPrimitiveID}}, bson.D{{"$push", bson.D{{"messages", message.ID.Hex()}}}})
	if err != nil {
		// delete the message if the update fails
		_, err = r.collectionMessage.DeleteOne(ctx, bson.D{{"_id", message.ID}})
		if err != nil {
			return nil, err
		}
		return nil, err

	}
	return message, nil
}

// GetMessages retrieves all messages from a room
func (r *messageRepository) GetMessages(ctx context.Context, roomID primitive.ObjectID) ([]*Message, error) {
	var roomRetrieved room.Room
	err := r.collectionRoom.FindOne(ctx, bson.D{{"_id", roomID}}).Decode(&roomRetrieved)
	if err != nil {
		fmt.Println("Error 1 : ", err)
		return nil, err

	}

	var messages []*Message
	for _, messageId := range roomRetrieved.Messages {
		var messageRetrieved *Message
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

	return messages, nil
}

func (r *messageRepository) GetMessage(ctx context.Context, messageID primitive.ObjectID) (*Message, error) {
	var message Message
	err := r.collectionMessage.FindOne(ctx, bson.D{{"_id", messageID}}).Decode(&message)
	if err != nil {
		return nil, err
	}
	return &message, nil

}

// DeleteMessage deletes a message from the database
func (r *messageRepository) DeleteMessage(ctx context.Context, messageID primitive.ObjectID) error {
	// Define the filter to find the message with the given ID
	filter := bson.M{"_id": messageID}

	// Delete the message from the messages collection
	_, err := r.collectionMessage.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// Update the rooms collection to remove the message ID from the list of messages
	update := bson.M{"$pull": bson.M{"messages": messageID.Hex()}}
	_, err = r.collectionRoom.UpdateMany(ctx, bson.M{}, update)
	if err != nil {
		// If the update fails, reinsert the message ID into the list of messages
		_, err := r.collectionRoom.UpdateMany(ctx, bson.M{}, bson.M{"$push": bson.M{"messages": messageID.Hex()}})
		if err != nil {
			return err
		}
		return err
	}

	return nil
}
