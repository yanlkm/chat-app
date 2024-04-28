package room

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room *Room) (*Room, error)
	CheckName(ctx context.Context, name string) error
	GetRoom(ctx context.Context, roomID primitive.ObjectID) (*Room, error)
	GetAllRooms(ctx context.Context) ([]Room, error)
	AddMember(ctx context.Context, roomID primitive.ObjectID, memberID primitive.ObjectID) (*Room, error)
	RemoveMember(ctx context.Context, roomID primitive.ObjectID, memberID primitive.ObjectID) (*Room, error)
	Delete(ctx context.Context, roomID primitive.ObjectID) error
}

type roomRepository struct {
	collection *mongo.Collection
}

func NewRoomRepository(collection *mongo.Collection) RoomRepository {
	return &roomRepository{collection: collection}
}

func (r *roomRepository) CreateRoom(ctx context.Context, room *Room) (*Room, error) {
	// add others required fields by default
	room = &Room{
		ID:          primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        room.Name,
		Description: room.Description,
		Creator:     room.Creator,
		Messages:    []string{},
		Hashtags:    []string{"#room"},
		Members:     []string{room.Creator},
	}
	_, err := r.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *roomRepository) CheckName(ctx context.Context, name string) error {
	var room Room
	err := r.collection.FindOne(ctx, bson.D{{"name", name}}).Decode(&room)
	if err != nil {
		return nil
	}
	return errors.New("Name already exists")

}
func (r *roomRepository) GetRoom(ctx context.Context, roomID primitive.ObjectID) (*Room, error) {
	var room Room
	err := r.collection.FindOne(ctx, bson.D{{"_id", roomID}}).Decode(&room)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) GetAllRooms(ctx context.Context) ([]Room, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var rooms []Room
	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil

}

func (r *roomRepository) AddMember(ctx context.Context, roomID primitive.ObjectID, memberID primitive.ObjectID) (*Room, error) {
	var room Room
	errCheck := r.collection.FindOne(ctx, bson.D{{"_id", roomID}, {"members", memberID.Hex()}}).Decode(&room)
	if errCheck == nil {
		return nil, errors.New("Member already added to room")
	}
	// debug
	fmt.Println("err : ", errCheck)
	_, err := r.collection.UpdateOne(ctx, bson.D{{"_id", roomID}}, bson.D{{"$push", bson.D{{"members", memberID.Hex()}}}})
	if err != nil {
		return nil, err
	}
	return r.GetRoom(ctx, roomID)
}

func (r *roomRepository) RemoveMember(ctx context.Context, roomID primitive.ObjectID, memberID primitive.ObjectID) (*Room, error) {
	// check if member exists
	var room Room
	errCheck := r.collection.FindOne(ctx, bson.D{{"_id", roomID}, {"members", memberID.Hex()}}).Decode(&room)
	if errCheck != nil {
		return nil, errors.New("Member already removed from room")
	}
	_, err := r.collection.UpdateOne(ctx, bson.D{{"_id", roomID}}, bson.D{{"$pull", bson.D{{"members", memberID.Hex()}}}})
	if err != nil {
		return nil, err
	}
	return r.GetRoom(ctx, roomID)
}

func (r *roomRepository) Delete(ctx context.Context, roomID primitive.ObjectID) error {

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": roomID})

	if err != nil || result.DeletedCount == 0 {
		return errors.New("Fail to delete room")
	}
	return nil
}
