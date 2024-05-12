package room

import (
	"chat-app/pkg/user"
	"context"
	"errors"
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
	collection      *mongo.Collection
	collectionUsers *mongo.Collection
}

func NewRoomRepository(collection *mongo.Collection, collectionUsers *mongo.Collection) RoomRepository {
	return &roomRepository{collection: collection, collectionUsers: collectionUsers}
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
	// check if member exists in users
	var member user.User
	errCheckUser := r.collectionUsers.FindOne(ctx, bson.D{{"_id", memberID}}).Decode(&member)

	if errCheckUser != nil {
		return nil, errors.New("Member does not exist")
	}
	var room Room
	// check if member exists in room
	errCheck := r.collection.FindOne(ctx, bson.D{{"_id", roomID}, {"members", memberID.Hex()}}).Decode(&room)
	if errCheck == nil {
		return nil, errors.New("Member already added to room")
	}
	// add room to rooms fields of user
	_, err := r.collectionUsers.UpdateOne(ctx,
		bson.D{{"_id", memberID}},
		bson.D{{"$push", bson.D{{"joinedRooms", roomID.Hex()}}}})
	if err != nil {
		return nil, err
	}
	_, err = r.collection.UpdateOne(ctx, bson.D{{"_id", roomID}}, bson.D{{"$push", bson.D{{"members", memberID.Hex()}}}})
	if err != nil {
		// delete last room from rooms fields of user
		_, err = r.collectionUsers.UpdateOne(ctx,
			bson.D{{"_id", memberID}},
			bson.D{{"$pull", bson.D{{"joinedRooms", roomID.Hex()}}}})
		return nil, err
	}
	return r.GetRoom(ctx, roomID)
}

func (r *roomRepository) RemoveMember(ctx context.Context, roomID primitive.ObjectID, memberID primitive.ObjectID) (*Room, error) {
	// check if member exists in users collection
	var member user.User
	errCheckUser := r.collectionUsers.FindOne(ctx, bson.D{{"_id", memberID}}).Decode(&member)
	if errCheckUser != nil {
		return nil, errors.New("Member does not exist")
	}
	// remove room from rooms fields of user
	_, err := r.collectionUsers.UpdateOne(ctx,
		bson.D{{"_id", memberID}},
		bson.D{{"$pull", bson.D{{"joinedRooms", roomID.Hex()}}}})
	if err != nil {
		return nil, err
	}
	var room Room
	// check if member exists in room
	errCheck := r.collection.FindOne(ctx, bson.D{{"_id", roomID}, {"members", memberID.Hex()}}).Decode(&room)
	if errCheck != nil {
		return nil, errors.New("Member already removed from room")
	}
	_, err = r.collection.UpdateOne(ctx, bson.D{{"_id", roomID}}, bson.D{{"$pull", bson.D{{"members", memberID.Hex()}}}})
	if err != nil {
		// add room to rooms fields of user
		_, err = r.collectionUsers.UpdateOne(ctx,
			bson.D{{"_id", memberID}},
			bson.D{{"$push", bson.D{{"joinedRooms", roomID.Hex()}}}})
		return nil, err
	}
	return r.GetRoom(ctx, roomID)
}

func (r *roomRepository) Delete(ctx context.Context, roomID primitive.ObjectID) error {
	// remove room from all users
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": roomID})
	if err != nil || result.DeletedCount == 0 {
		return errors.New("Fail to delete room")
	}
	return nil
}
