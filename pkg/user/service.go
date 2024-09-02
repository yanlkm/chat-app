package user

import (
	"context"
	"time"
)

// UserService provides user operations.
type UserService interface {
	CreateUser(ctx context.Context, user *UserEntity) error
	CheckEmail(ctx context.Context, email string) error
	CheckUsername(ctx context.Context, username string) error
	GetUser(ctx context.Context, id string) (*UserEntity, error)
	GetAllUsers(ctx context.Context) ([]UserEntity, error)
	UpdateUser(ctx context.Context, id string, username string) error
	UpdatePassword(ctx context.Context, id string, newPassword string) error
	BanUser(ctx context.Context, idBanner string, idBanned string) error
	UnBanUser(ctx context.Context, idBanner string, idBanned string) error
	DeleteUser(ctx context.Context, id string) error
}

// userService implements UserService.
type userService struct {
	repo UserRepository
}

// NewUserService creates a new user service.
func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// CreateUser creates a new user.
func (s *userService) CreateUser(ctx context.Context, user *UserEntity) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return s.repo.Create(ctx, user)
}

// CheckEmail checks if an email is already in use.
func (s *userService) CheckEmail(ctx context.Context, email string) error {
	return s.repo.CheckEmail(ctx, email)
}

// CheckUsername checks if a username is already in use.
func (s *userService) CheckUsername(ctx context.Context, username string) error {
	return s.repo.CheckUsername(ctx, username)
}

// GetUser retrieves a user by ID.
func (s *userService) GetUser(ctx context.Context, id string) (*UserEntity, error) {
	return s.repo.Read(ctx, id)
}

// GetAllUsers retrieves all users.
func (s *userService) GetAllUsers(ctx context.Context) ([]UserEntity, error) {
	return s.repo.ReadUsers(ctx)
}

// UpdateUser updates a user's username.
func (s *userService) UpdateUser(ctx context.Context, id string, username string) error {
	return s.repo.Update(ctx, id, username)
}

// UpdatePassword updates a user's password.
func (s *userService) UpdatePassword(ctx context.Context, id string, newPassword string) error {
	return s.repo.UpdatePassword(ctx, id, newPassword)
}

// BanUser bans a user.
func (s *userService) BanUser(ctx context.Context, idBanner string, idBanned string) error {
	return s.repo.BanUser(ctx, idBanner, idBanned)
}

// UnBanUser unbans a user.
func (s *userService) UnBanUser(ctx context.Context, idBanner string, idBanned string) error {
	return s.repo.UnBanUser(ctx, idBanner, idBanned)
}

// DeleteUser deletes a user.
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
