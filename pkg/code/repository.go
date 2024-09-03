package code

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// repository to create and update code
type CodeRepository interface {
	Create(ctx context.Context, code *CodeEntity) error
	Update(ctx context.Context, codeString *string) error
	Check(ctx context.Context, codeString *string) (bool, error)
}

type codeRepository struct {
	collection *mongo.Collection
}

func NewCodeRepository(collection *mongo.Collection) CodeRepository {
	return &codeRepository{collection: collection}
}

// Create a new code
func (r *codeRepository) Create(ctx context.Context, code *CodeEntity) error {
	// add the current time to the code and set is_used to false
	codeModel := &CodeModel{
		Code:      code.Code,
		CreatedAt: time.Now(),
		IsUsed:    false,
	}
	// Check if the code is unique in the database
	filter := bson.M{"code": codeModel.Code}
	err := r.collection.FindOne(ctx, filter).Err()
	if err == nil {
		// The code already exists in the database
		return errors.New("code already exists")
	}
	if err != mongo.ErrNoDocuments {
		// An unexpected error occurred while searching for the code
		return err
	}

	// Insert the code into the database
	_, err = r.collection.InsertOne(ctx, codeModel)
	if err != nil {
		// An error occurred while inserting the code
		return err
	}

	return nil
}

// Update an existing code
func (r *codeRepository) Update(ctx context.Context, codeString *string) error {
	// Update the code status to true when it's used
	filter := bson.M{"code": *codeString}
	update := bson.M{"$set": bson.M{"isUsed": true}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		// An error occurred while updating the code status
		return err
	}

	return nil
}

// Check if a code existing
func (r *codeRepository) Check(ctx context.Context, codeString *string) (bool, error) {
	// Check if the code exists in the database and is not used
	filter := bson.M{"code": *codeString, "isUsed": false}
	err := r.collection.FindOne(ctx, filter).Err()
	if err == mongo.ErrNoDocuments {
		// The code does not exist in the database
		return false, nil
	}
	if err != nil {
		// An unexpected error occurred while searching for the code
		return false, err
	}

	return true, nil
}
