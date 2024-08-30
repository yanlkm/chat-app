package user

import "time"

type UserModel struct {
	ID           string    `bson:"_id,omitempty"`
	Username     string    `bson:"username,omitempty"`
	Email        string    `bson:"email,omitempty"`
	Password     string    `bson:"password,omitempty"`
	CreatedAt    time.Time `bson:"createdAt,omitempty"`
	UpdatedAt    time.Time `bson:"updatedAt,omitempty"`
	Role         string    `bson:"role,omitempty"`
	Validity     string    `bson:"validity,omitempty"`
	JoinedSalons []string  `bson:"joinedRooms,omitempty"`
}

func ModelToEntity(model *UserModel) *UserEntity {
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

func EntityToModel(entity *UserEntity) *UserModel {
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
