package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	CheckEmail(ctx context.Context, email string) error
	CheckUsername(ctx context.Context, username string) error
	GetUser(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, username string) error
	UpdatePassword(ctx context.Context, id primitive.ObjectID, newPassword string) error
	BanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error
	UnBanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error
	DeleteUser(ctx context.Context, id primitive.ObjectID) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return s.repo.Create(ctx, user)
}

func (s *userService) CheckEmail(ctx context.Context, email string) error {
	return s.repo.CheckEmail(ctx, email)
}
func (s *userService) CheckUsername(ctx context.Context, username string) error {
	return s.repo.CheckUsername(ctx, username)
}

func (s *userService) GetUser(ctx context.Context, id primitive.ObjectID) (*User, error) {
	return s.repo.Read(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]User, error) {
	return s.repo.ReadUsers(ctx)
}

func (s *userService) UpdateUser(ctx context.Context, id primitive.ObjectID, username string) error {
	return s.repo.Update(ctx, id, username)
}

func (s *userService) UpdatePassword(ctx context.Context, id primitive.ObjectID, newPassword string) error {
	return s.repo.UpdatePassword(ctx, id, newPassword)
}

func (s *userService) BanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error {
	return s.repo.BanUser(ctx, idBanner, idBanned)
}

func (s *userService) UnBanUser(ctx context.Context, idBanner primitive.ObjectID, idBanned primitive.ObjectID) error {
	return s.repo.UnBanUser(ctx, idBanner, idBanned)
}

func (s *userService) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}
