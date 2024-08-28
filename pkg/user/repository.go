package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Read(ctx context.Context, id primitive.ObjectID) (*User, error)
	ReadUsers(ctx context.Context) ([]User, error)
	CheckEmail(ctx context.Context, email string) error
	CheckUsername(ctx context.Context, username string) error
	Update(ctx context.Context, id primitive.ObjectID, username string) error
	UpdatePassword(ctx context.Context, id primitive.ObjectID, newPassword string) error
	BanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error
	UnBanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// userRepository represents the repository for managing users.
type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository.
func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{collection: collection}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// CheckUsername checks if the username already exists in the database.
func (r *userRepository) CheckUsername(ctx context.Context, username string) error {
	var user User
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
	var user User
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

// Read returns the user with the provided ID.
func (r *userRepository) Read(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// get all users except admins and return them
func (r *userRepository) ReadUsers(ctx context.Context) ([]User, error) {
	var users []User
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		if user.Role != "admin" {
			users = append(users, user)
		}
	}
	return users, nil
}

// Update updates the username for a user.
func (r *userRepository) Update(ctx context.Context, id primitive.ObjectID, username string) error {
	// check if user exists
	var user User
	// check if username is unique
	err := r.collection.FindOne(ctx, bson.D{{"username", username}}).Decode(&user)
	// if the username already exists and it is not the user's username by id converted to string
	if err == nil && user.ID != id.Hex() {
		return errors.New("Username already exists")
	}
	// Update the username in the database
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"username": username, "updatedAt": time.Now()}})
	if err != nil {
		return err
	}
	return nil
}

// UpdatePassword updates the password for a user.
func (r *userRepository) UpdatePassword(ctx context.Context, id primitive.ObjectID, newPassword string) error {
	// update the password in the database
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"password": newPassword, "updatedAt": time.Now()}})
	if err != nil {
		return err
	}
	return nil
}

// BanUser bans a user from the platform.
func (r *userRepository) BanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error {
	// Check if banner exists and is a valid admin
	var banner User
	err := r.collection.FindOne(ctx, bson.M{"_id": idBanner}).Decode(&banner)
	if err != nil {
		return errors.New("Error banning user")
	}
	if banner.Role != "admin" || banner.Validity != "valid" {
		return errors.New("Error banning user")
	}
	// Check if banned user exists
	var banned User
	err = r.collection.FindOne(ctx, bson.M{"_id": idBanned}).Decode(&banned)
	if err != nil {
		return errors.New("Error banning user")
	}
	// Check if banned user is not an admin
	if banned.Role == "admin" {
		return errors.New("Error banning user")
	}

	// Ban the user
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": idBanned}, bson.M{"$set": bson.M{"validity": "invalid"}})
	if err != nil {
		return err
	}
	return nil
}

// UnBanUser unbans a user from the platform.
func (r *userRepository) UnBanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error {
	// Check if banner exists and is a valid admin
	var banner User
	err := r.collection.FindOne(ctx, bson.M{"_id": idBanner}).Decode(&banner)
	if err != nil {
		return errors.New("Error unbanning user")
	}
	if banner.Role != "admin" || banner.Validity != "valid" {
		return errors.New("Error unbanning user")
	}
	// Check if banned user exists
	var banned User
	err = r.collection.FindOne(ctx, bson.M{"_id": idBanned}).Decode(&banned)
	if err != nil {
		return errors.New("Error unbanning user")
	}
	// Unban the user
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": idBanned}, bson.M{"$set": bson.M{"validity": "valid"}})
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a user from the platform.
func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
