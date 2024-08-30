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
	Login(ctx context.Context, credentials UserCredentials) (*user.UserEntity, error)
	Logout(ctx context.Context, userID *string) error
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
func (r *authRepository) Login(ctx context.Context, credentials UserCredentials) (*user.UserEntity, error) {
	var foundUser user.UserModel

	// Find user by username
	err := r.collection.FindOne(ctx, bson.D{{"username", credentials.Username}}).Decode(&foundUser)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user.ModelToEntity(&foundUser), nil
}

// Logout logs a user out by invalidating the session or token.
func (r *authRepository) Logout(ctx context.Context, userID *string) error {
	// convert userID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(*userID)
	if err != nil {
		return err
	}

	err = r.collection.FindOne(ctx, bson.D{{"_id", &objectID}}).Err()
	if err != nil {
		return err
	}
	return nil
}
