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

// RoomRepository is the interface that wraps the basic room repository methods.
type RoomRepository interface {
	CreateRoom(ctx context.Context, room *RoomEntity) (*RoomEntity, error)
	CheckName(ctx context.Context, name string) error
	GetRoom(ctx context.Context, roomID string) (*RoomEntity, error)
	GetUserRooms(ctx context.Context, userID string) ([]RoomEntity, error)
	GetAllRooms(ctx context.Context) ([]RoomEntity, error)
	GetRoomsCreatedByAdmin(ctx context.Context, adminID string) ([]RoomEntity, error)
	AddMember(ctx context.Context, roomID string, memberID string) (*RoomEntity, error)
	RemoveMember(ctx context.Context, roomID string, memberID string) (*RoomEntity, error)
	AddHashtag(ctx context.Context, roomID string, hashtag string) (*RoomEntity, error)
	RemoveHashtag(ctx context.Context, roomID string, hashtag string) (*RoomEntity, error)
	Delete(ctx context.Context, roomID string) error
}

// roomRepository is the implementation of the RoomRepository interface.
type roomRepository struct {
	collection      *mongo.Collection
	collectionUsers *mongo.Collection
}

// NewRoomRepository creates a new room repository.
func NewRoomRepository(collection *mongo.Collection, collectionUsers *mongo.Collection) RoomRepository {
	return &roomRepository{collection: collection, collectionUsers: collectionUsers}
}

// CreateRoom creates a new room in the database.
func (r *roomRepository) CreateRoom(ctx context.Context, room *RoomEntity) (*RoomEntity, error) {
	// check if creator is a valid objectID
	_, err := primitive.ObjectIDFromHex(room.Creator)
	if err != nil {
		return nil, errors.New(" Invalid creator")
	}

	// add others required fields by default
	roomModel := &RoomModel{
		ID:          primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Name:        room.Name,
		Description: room.Description,
		Creator:     room.Creator,
		Messages:    []string{},
		Hashtags:    []string{"#room"}, // add default hashtag to room
		Members:     []string{room.Creator},
	}
	// check if creator exists in users
	var userCheck user.UserModel
	// convert creator to objectID
	roomCreatorObjectID, err := primitive.ObjectIDFromHex(roomModel.Creator)
	if err != nil {
		return nil, errors.New(" The creator does not exist")
	}

	// check if the user room creator exists
	errCheckUser := r.collectionUsers.FindOne(ctx, bson.D{{"_id", roomCreatorObjectID}}).Decode(&userCheck)
	if errCheckUser != nil {
		return nil, errors.New(" The creator does not exist")
	}

	// check if room creator is valid and admin
	if userCheck.Validity != "valid" || userCheck.Role != "admin" {
		return nil, errors.New(" The creator is not a valid user or an admin")
	}
	// insert room model to database
	_, err = r.collection.InsertOne(ctx, roomModel)
	if err != nil {
		return nil, err
	}
	return room, nil
}

// CheckName checks if the name already exists in the database.
func (r *roomRepository) CheckName(ctx context.Context, name string) error {
	// invoke Room Model
	var room RoomModel
	err := r.collection.FindOne(ctx, bson.D{{"name", name}}).Decode(&room)
	if err != nil {
		return nil
	}
	return errors.New(" Name already exists")

}

// GetRoom returns a room by its ID
func (r *roomRepository) GetRoom(ctx context.Context, roomID string) (*RoomEntity, error) {
	// convert roomID to objectID
	roomIDObjectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, errors.New(" Invalid room ID")
	}

	// invoke Room Model
	var room RoomModel
	err = r.collection.FindOne(ctx, bson.D{{"_id", roomIDObjectID}}).Decode(&room)
	if err != nil {
		return nil, err
	}
	// return room model mapped into entity
	return ModelToEntity(&room), nil
}

// GetUserRooms returns all rooms where the user is a member
func (r *roomRepository) GetUserRooms(ctx context.Context, userID string) ([]RoomEntity, error) {
	// convert userID to objectID
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New(" Invalid user ID")
	}
	// check if userCheck exists in users
	var userCheck user.UserModel
	errCheckUser := r.collectionUsers.FindOne(ctx, bson.D{{"_id", userObjectID}}).Decode(&userCheck)
	if errCheckUser != nil {
		return nil, errors.New(" User does not exist")
	}

	// define a cursor of members
	cursor, err := r.collection.Find(ctx, bson.D{{"members", userID}})
	if err != nil {
		return nil, err
	}
	// defer cursor
	defer cursor.Close(ctx)
	// fin rooms of user
	var rooms []RoomModel
	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}
	// map  rooms model into entities
	// convert list of rooms to list of room entities
	roomsEntities := make([]RoomEntity, 0)
	for _, room := range rooms {
		roomsEntities = append(roomsEntities, *ModelToEntity(&room))
	}
	// return entities
	return roomsEntities, nil
}

// GetAllRooms returns all rooms
func (r *roomRepository) GetAllRooms(ctx context.Context) ([]RoomEntity, error) {
	// get all rooms
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var rooms []RoomModel
	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}
	// convert list of rooms to list of room entities
	roomsEntities := make([]RoomEntity, 0)
	for _, room := range rooms {
		roomsEntities = append(roomsEntities, *ModelToEntity(&room))
	}
	return roomsEntities, nil

}

func (r *roomRepository) GetRoomsCreatedByAdmin(ctx context.Context, adminID string) ([]RoomEntity, error) {
	// convert adminID to objectID
	adminIDObjectID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return nil, errors.New(" Invalid admin ID")
	}
	// check if userCheck exists in users
	var adminCheck user.UserModel
	errCheckAdmin := r.collectionUsers.FindOne(ctx, bson.D{{"_id", adminIDObjectID}}).Decode(&adminCheck)
	if errCheckAdmin != nil {
		return nil, errors.New("User does not exist")
	}
	cursor, err := r.collection.Find(ctx, bson.D{{"creator", adminID}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var rooms []RoomModel
	if err = cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}
	// convert list of rooms to list of room entities
	roomsEntities := make([]RoomEntity, 0)
	for _, room := range rooms {
		roomsEntities = append(roomsEntities, *ModelToEntity(&room))
	}
	return roomsEntities, nil

}

func (r *roomRepository) AddMember(ctx context.Context, roomID string, memberID string) (*RoomEntity, error) {
	// check IDs by converting them to objectIDs
	roomIDObjectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, errors.New(" Invalid room ID")
	}
	memberIDObjectID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, errors.New(" Invalid member ID")
	}

	// check if member exists in users
	var member user.UserModel
	errCheckUser := r.collectionUsers.FindOne(ctx, bson.D{{"_id", memberIDObjectID}}).Decode(&member)

	if errCheckUser != nil {
		return nil, errors.New(" Member does not exist")
	}
	var room RoomModel
	// check if member exists in room
	errCheck := r.collection.FindOne(ctx, bson.D{{"_id", roomIDObjectID}, {"members", memberID}}).Decode(&room)
	if errCheck == nil {
		return nil, errors.New(" Member already added to room")
	}
	// add room to rooms fields of user
	_, err = r.collectionUsers.UpdateOne(ctx,
		bson.D{{"_id", memberIDObjectID}},
		bson.D{{"$push", bson.D{{"joinedRooms", roomID}}}})
	if err != nil {
		return nil, err
	}
	_, err = r.collection.UpdateOne(ctx, bson.D{{"_id", roomIDObjectID}}, bson.D{{"$push", bson.D{{"members", memberID}}}})
	if err != nil {
		// delete last room from rooms fields of user
		_, err = r.collectionUsers.UpdateOne(ctx,
			bson.D{{"_id", memberIDObjectID}},
			bson.D{{"$pull", bson.D{{"joinedRooms", roomID}}}})
		return nil, err
	}
	return r.GetRoom(ctx, roomID)
}

func (r *roomRepository) RemoveMember(ctx context.Context, roomID string, memberID string) (*RoomEntity, error) {
	// check IDs by converting them to objectIDs
	roomIDObjectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, errors.New(" Invalid room ID")
	}
	memberIDObjectID, err := primitive.ObjectIDFromHex(memberID)
	if err != nil {
		return nil, errors.New(" Invalid member ID")
	}

	// check if member exists in users collection
	var member user.UserModel
	errCheckUser := r.collectionUsers.FindOne(ctx, bson.D{{"_id", memberIDObjectID}}).Decode(&member)
	if errCheckUser != nil {
		return nil, errors.New(" Member does not exist")
	}
	// check if member is the creator of the room
	var room RoomModel
	errCheck := r.collection.FindOne(ctx, bson.D{{"_id", roomIDObjectID}, {"creator", memberID}}).Decode(&room)
	if errCheck == nil {
		return nil, errors.New(" Member is the creator of the room")
	}
	// remove room from rooms fields of user
	_, err = r.collectionUsers.UpdateOne(ctx,
		bson.D{{"_id", memberIDObjectID}},
		bson.D{{"$pull", bson.D{{"joinedRooms", roomID}}}})
	if err != nil {
		return nil, err
	}

	// check if member exists in room
	errCheck = r.collection.FindOne(ctx, bson.D{{"_id", roomIDObjectID}, {"members", memberID}}).Decode(&room)
	if errCheck != nil {
		return nil, errors.New(" Member already removed from room")
	}
	_, err = r.collection.UpdateOne(ctx, bson.D{{"_id", roomIDObjectID}}, bson.D{{"$pull", bson.D{{"members", memberID}}}})
	if err != nil {
		// add room to rooms fields of user
		_, err = r.collectionUsers.UpdateOne(ctx,
			bson.D{{"_id", memberIDObjectID}},
			bson.D{{"$push", bson.D{{"joinedRooms", roomID}}}})
		return nil, err
	}
	return r.GetRoom(ctx, roomID)
}

func (r *roomRepository) AddHashtag(ctx context.Context, roomID string, hashtag string) (*RoomEntity, error) {

	// convert roomID to objectID
	roomIDObjectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, errors.New(" Invalid room ID")
	}

	//check if room exists
	var room RoomModel
	errCheck := r.collection.FindOne(ctx, bson.M{"_id": roomIDObjectID}).Decode(&room)
	if errCheck != nil {
		return nil, errors.New(" Room does not exist")
	}

	//check if hashtag already exists in hashtag array
	errHashtag := r.collection.FindOne(ctx, bson.D{{"_id", roomIDObjectID}, {"hashtags", hashtag}})
	if errHashtag == nil {
		return nil, errors.New(" Hashtag already added to room")
	}

	// add hashtag to room
	_, err = r.collection.UpdateOne(ctx, bson.D{{"_id", roomIDObjectID}}, bson.D{{"$push", bson.D{{"hashtags", hashtag}}}})
	if err != nil {
		return nil, err
	}

	// update last update field
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": roomIDObjectID}, bson.M{"$set": bson.M{"updatedAt": time.Now()}})
	if err != nil {

		return nil, err
	}

	return r.GetRoom(ctx, roomID)
}

func (r *roomRepository) RemoveHashtag(ctx context.Context, roomID string, hashtag string) (*RoomEntity, error) {

	// convert roomID to objectID
	roomIDObjectID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, errors.New(" Invalid room ID")
	}

	//check if room exists
	var room RoomModel
	errCheck := r.collection.FindOne(ctx, bson.M{"_id": roomIDObjectID}).Decode(&room)
	if errCheck != nil {
		return nil, errors.New(" Room does not exist")
	}
	// check if hashtag doesn't exist in hashtag array
	errHashtag := r.collection.FindOne(ctx, bson.D{{"_id", roomIDObjectID}, {"hashtags", hashtag}})
	if errHashtag == nil {
		return nil, errors.New(" Hashtag already removed from room")
	}

	// check if hashtag array is not empty or there is at least two hashtags
	if len(room.Hashtags) < 2 {
		return nil, errors.New(" Hashtag array is empty or there is only one hashtag")
	}

	_, err = r.collection.UpdateOne(ctx, bson.D{{"_id", roomIDObjectID}}, bson.D{{"$pull", bson.D{{"hashtags", hashtag}}}})
	if err != nil {
		return nil, err
	}

	// update last update field
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": roomIDObjectID}, bson.M{"$set": bson.M{"updatedAt": time.Now()}})
	if err != nil {

		return nil, err
	}

	return r.GetRoom(ctx, roomID)
}

func (r *roomRepository) Delete(ctx context.Context, roomID string) error {
	// remove room from all users
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": roomID})
	if err != nil || result.DeletedCount == 0 {
		return errors.New(" Fail to delete room")
	}
	return nil
}
