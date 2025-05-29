package service

import (
	"context"
	"fmt"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/store"
)

type AuthService struct {
	store store.Storage
}

func NewAuthService(store store.Storage) *AuthService {
	return &AuthService{store: store}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterUserRequest) (*models.User, error) {

	hashedPassword, err := models.NewPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = s.store.User.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	
	// user, err := s.store.User.GetByEmail(ctx, email)
	// if err != nil {
	//     return nil, err
	// }

	// // Check if the password matches
	// if err := user.Password.Compare(password); err != nil {
	//     return nil, fmt.Errorf("invalid credentials")
	// }

	// return user, nil

	return nil, fmt.Errorf("login method needs GetByEmail implementation in UserRepository")
}
