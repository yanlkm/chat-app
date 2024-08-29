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
	Create(ctx context.Context, user *UserEntity) error
	Read(ctx context.Context, id primitive.ObjectID) (*UserEntity, error)
	ReadUsers(ctx context.Context) ([]UserEntity, error)
	CheckEmail(ctx context.Context, email string) error
	CheckUsername(ctx context.Context, username string) error
	Update(ctx context.Context, id primitive.ObjectID, username string) error
	UpdatePassword(ctx context.Context, id primitive.ObjectID, newPassword string) error
	BanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error
	UnBanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{collection: collection}
}

func modelToEntity(model *UserModel) *UserEntity {
	return &UserEntity{
		ID:           model.ID,
		Username:     model.Username,
		Email:        model.Email,
		Password:     model.Password,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
		Role:         model.Role,
		Validity:     model.Validity,
		JoinedSalons: model.JoinedSalons,
	}
}

func entityToModel(entity *UserEntity) *UserModel {
	return &UserModel{
		ID:           entity.ID,
		Username:     entity.Username,
		Email:        entity.Email,
		Password:     entity.Password,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
		Role:         entity.Role,
		Validity:     entity.Validity,
		JoinedSalons: entity.JoinedSalons,
	}
}

// Create creates a new user in the database.
func (r *userRepository) Create(ctx context.Context, user *UserEntity) error {
	model := entityToModel(user)
	_, err := r.collection.InsertOne(ctx, model)
	return err
}

// Read returns the user with the provided ID.
func (r *userRepository) Read(ctx context.Context, id primitive.ObjectID) (*UserEntity, error) {
	var model UserModel
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&model)
	if err != nil {
		return nil, err
	}
	return modelToEntity(&model), nil
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

// get all users except admins and return them
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
			users = append(users, *modelToEntity(&user))
		}
	}
	return users, nil
}

// Update updates the username for a user.
func (r *userRepository) Update(ctx context.Context, id primitive.ObjectID, username string) error {
	// check if user exists
	var user UserModel
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
	var banner UserModel
	err := r.collection.FindOne(ctx, bson.M{"_id": idBanner}).Decode(&banner)
	if err != nil {
		return errors.New("Error banning user")
	}
	if banner.Role != "admin" || banner.Validity != "valid" {
		return errors.New("Error banning user")
	}
	// Check if banned user exists
	var banned UserModel
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
	var banner UserModel
	err := r.collection.FindOne(ctx, bson.M{"_id": idBanner}).Decode(&banner)
	if err != nil {
		return errors.New("Error unbanning user")
	}
	if banner.Role != "admin" || banner.Validity != "valid" {
		return errors.New("Error unbanning user")
	}
	// Check if banned user exists
	var banned UserModel
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
