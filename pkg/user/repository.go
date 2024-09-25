package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// UserRepository is the interface that wraps the basic CRUD operations for the user entity.
type UserRepository interface {
	Create(ctx context.Context, user *UserEntity) error
	Read(ctx context.Context, id string) (*UserEntity, error)
	ReadUsers(ctx context.Context) ([]UserEntity, error)
	CheckEmail(ctx context.Context, email string) error
	CheckUsername(ctx context.Context, username string) error
	Update(ctx context.Context, id string, username string) error
	UpdatePassword(ctx context.Context, id string, newPassword string) error
	BanUser(ctx context.Context, idBanner string, idBanned string) error
	UnBanUser(ctx context.Context, idBanner string, idBanned string) error
	Delete(ctx context.Context, id string) error
}

// userRepository represents the repository for the user entity.
type userRepository struct {
	collection        *mongo.Collection
	collectionMessage *mongo.Collection
}

// NewUserRepository creates a new user repository.
func NewUserRepository(collection *mongo.Collection, collectionMessage *mongo.Collection) UserRepository {
	return &userRepository{collection: collection, collectionMessage: collectionMessage}
}

// Create creates a new user in the database.
func (r *userRepository) Create(ctx context.Context, user *UserEntity) error {
	model := EntityToModel(user)
	_, err := r.collection.InsertOne(ctx, model)
	return err
}

// Read returns the user with the provided ID.
func (r *userRepository) Read(ctx context.Context, id string) (*UserEntity, error) {
	// convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var model UserModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&model)
	if err != nil {
		return nil, err
	}
	return ModelToEntity(&model), nil
}

// CheckUsername checks if the username already exists in the database.
func (r *userRepository) CheckUsername(ctx context.Context, username string) error {
	var user UserModel
	err := r.collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			return errors.New("username already exists")
		}
	}
	return errors.New("username already exists")
}

// CheckEmail checks if the email already exists in the database.
func (r *userRepository) CheckEmail(ctx context.Context, email string) error {
	var user UserModel
	err := r.collection.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			return errors.New("email already exists")
		}
	}
	return errors.New("email already exists")
}

// ReadUsers get all users except admins and return them
func (r *userRepository) ReadUsers(ctx context.Context) ([]UserEntity, error) {
	var users []UserEntity
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user UserModel
		cursor.Decode(&user)
		if user.Role != "admin" {
			users = append(users, *ModelToEntity(&user))
		}
	}
	return users, nil
}

// Update updates the username for a user.
func (r *userRepository) Update(ctx context.Context, id string, username string) error {
	// convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// check if user exists
	var user UserModel
	// check if username is unique
	err = r.collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&user)
	// if the username already exists and it is not the user's username by id converted to string
	if err == nil && user.ID != id {
		return errors.New("Username already exists")
	}
	// store the old username if error occurs
	var oldUsername string
	oldUsername = user.Username

	// Update the username in the database
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"username": username, "updatedAt": time.Now()}})
	if err != nil {
		return err
	}

	// Update the username in the messages collection
	// Define the filter to find the collection message with the given userID in key "userId"
	filter := bson.M{"userId": id}

	// Define the update to change the username to the newUsername
	update := bson.M{"$set": bson.M{"username": username}}
	// find the messages with the filter and update the username
	_, err = r.collectionMessage.UpdateMany(ctx, filter, update)
	// if error occurs
	if err != nil {
		// if error occurs, revert the username to the old username
		_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"username": oldUsername, "updatedAt": time.Now()}})
		if err != nil {
			return err
		}
		return err
	}
	// return nil if the update is successful
	return nil
}

// UpdatePassword updates the password for a user.
func (r *userRepository) UpdatePassword(ctx context.Context, id string, newPassword string) error {
	// convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// update the password in the database
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": bson.M{"password": newPassword, "updatedAt": time.Now()}})
	if err != nil {
		return err
	}
	return nil
}

// BanUser bans a user from the platform.
func (r *userRepository) BanUser(ctx context.Context, idBanner string, idBanned string) error {
	// convert ids to ObjectIDs
	objectIDBanner, err := primitive.ObjectIDFromHex(idBanner)
	if err != nil {
		return err
	}
	objectIDBanned, err := primitive.ObjectIDFromHex(idBanned)
	if err != nil {
		return err
	}
	// Check if banner exists and is a valid admin
	var banner UserModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objectIDBanner}).Decode(&banner)
	if err != nil {
		return errors.New("Error banning user")
	}
	if banner.Role != "admin" || banner.Validity != "valid" {
		return errors.New("Error banning user")
	}
	// Check if banned user exists
	var banned UserModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objectIDBanned}).Decode(&banned)
	if err != nil {
		return errors.New("Error banning user")
	}
	// Check if banned user is not an admin
	if banned.Role == "admin" {
		return errors.New("Error banning user")
	}

	// Ban the user
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectIDBanned}, bson.M{"$set": bson.M{"validity": "invalid"}})
	if err != nil {
		return err
	}
	return nil
}

// UnBanUser unbans a user from the platform.
func (r *userRepository) UnBanUser(ctx context.Context, idBanner string, idBanned string) error {
	// convert ids to ObjectIDs
	objectIDBanner, err := primitive.ObjectIDFromHex(idBanner)
	if err != nil {
		return err
	}
	objectIDBanned, err := primitive.ObjectIDFromHex(idBanned)
	if err != nil {
		return err
	}

	// Check if banner exists and is a valid admin
	var banner UserModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objectIDBanner}).Decode(&banner)
	if err != nil {
		return errors.New("Error unbanning user")
	}
	if banner.Role != "admin" || banner.Validity != "valid" {
		return errors.New("Error unbanning user")
	}
	// Check if banned user exists
	var banned UserModel
	err = r.collection.FindOne(ctx, bson.M{"_id": objectIDBanned}).Decode(&banned)
	if err != nil {
		return errors.New("Error unbanning user")
	}
	// Unban the user
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectIDBanned}, bson.M{"$set": bson.M{"validity": "valid"}})
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a user from the platform.
func (r *userRepository) Delete(ctx context.Context, id string) error {
	// convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	return nil
}
