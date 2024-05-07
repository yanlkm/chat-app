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
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type userRepository struct {
	collection *mongo.Collection
}

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

func (r *userRepository) Read(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

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
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, id primitive.ObjectID, username string) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"username": username, "updatedAt": time.Now()}})
	if err != nil {
		return err
	}
	return nil
}

// UpdatePassword updates the password for a user.
func (r *userRepository) UpdatePassword(ctx context.Context, id primitive.ObjectID, newPassword string) error {

	// Update the password in the database
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"password": newPassword, "updatedAt": time.Now()}})
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}
