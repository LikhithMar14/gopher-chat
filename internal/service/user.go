package service

import (
	"context"
	"fmt"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
)

type UserService struct {
	store store.Storage
}

func NewUserService(store store.Storage) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.store.User.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
	if err := s.validateCreateUserRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	if err := s.store.User.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *UserService) validateCreateUserRequest(req models.CreateUserRequest) error {
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

func (s *UserService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.store.User.GetByID(ctx, userID)
	if err != nil {
		switch {
		case err == store.ErrNotFound:
			return nil, apperrors.ErrUserNotFound
		default:
			return nil, fmt.Errorf("failed to get user by id: %w", err)
		}
	}
	return user, nil
}

func (s *UserService) GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(utils.UserIDKey).(*models.User)
	return user, ok
}