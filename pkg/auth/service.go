package auth

import (
	"chat-app/pkg/user"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthService interface {
	LoginUser(ctx context.Context, userLogin *UserCredentials) (*user.User, error)
	LogoutUser(ctx context.Context, id *primitive.ObjectID) error
}

type authService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) LoginUser(ctx context.Context, userLogin *UserCredentials) (*user.User, error) {
	return s.repo.Login(ctx, *userLogin)
}
func (s *authService) LogoutUser(ctx context.Context, id *primitive.ObjectID) error {
	return s.repo.Logout(ctx, id)
}
