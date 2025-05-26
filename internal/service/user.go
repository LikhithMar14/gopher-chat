package service

import (
	"context"
	"fmt"

	"github.com/LikhithMar14/gopher-chat/internal/store"
)

type UserService struct {
	store store.Storage
}

func NewUserService(store store.Storage) *UserService {
	return &UserService{
		store: store,
	}
}


func (s *UserService) GetUsers(ctx context.Context) ([]store.User, error) {
	users, err := s.store.User.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (s *UserService) CreateUser(ctx context.Context, req CreateUserRequest) (*store.User, error) {
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	user := &store.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	if err := s.store.User.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *UserService) validateCreateUserRequest(req CreateUserRequest) error {
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}
