package auth

import (
	"chat-app/pkg/user"
	"context"
)

type AuthService interface {
	LoginUser(ctx context.Context, userLogin *UserCredentials) (*user.UserEntity, error)
	LogoutUser(ctx context.Context, id *string) error
}

type authService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) LoginUser(ctx context.Context, userLogin *UserCredentials) (*user.UserEntity, error) {
	return s.repo.Login(ctx, *userLogin)
}
func (s *authService) LogoutUser(ctx context.Context, id *string) error {
	return s.repo.Logout(ctx, id)
}
