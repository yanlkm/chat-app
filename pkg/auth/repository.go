package auth

import (
	"chat-app/pkg/user"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuthRepository defines the interface for authentication repository operations.
type AuthRepository interface {
	Login(ctx context.Context, credentials UserCredentials) (*user.User, error)
	Logout(ctx context.Context, userID *primitive.ObjectID) error
}

// authRepository is the concrete implementation of AuthRepository.
type authRepository struct {
	collection *mongo.Collection
}

// NewAuthRepository creates a new instance of AuthRepository.
func NewAuthRepository(collection *mongo.Collection) AuthRepository {
	return &authRepository{collection: collection}
}

// Login attempts to authenticate a user with the provided credentials.
func (r *authRepository) Login(ctx context.Context, credentials UserCredentials) (*user.User, error) {
	var foundUser user.User

	// Find user by username
	err := r.collection.FindOne(ctx, bson.D{{"username", credentials.Username}}).Decode(&foundUser)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &foundUser, nil
}

// Logout logs a user out by invalidating the session or token.
func (r *authRepository) Logout(ctx context.Context, userID *primitive.ObjectID) error {
	err := r.collection.FindOne(ctx, bson.D{{"_id", &userID}}).Err()
	if err != nil {
		return err
	}
	return nil
}
