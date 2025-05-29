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

func (s *UserService) FollowUser(ctx context.Context, followerID int64, userID int64) error {
	_, err := s.store.User.GetByID(ctx, userID)
	if err != nil {
		switch {
		case err == store.ErrNotFound:
			return apperrors.ErrUserNotFound
		default:
			return fmt.Errorf("failed to verify target user: %w", err)
		}
	}

	if err := s.store.User.FollowUser(ctx, userID, followerID); err != nil {
		return fmt.Errorf("failed to follow user: %w", err)
	}
	return nil
}

func (s *UserService) UnfollowUser(ctx context.Context, followerID int64, userID int64) error {
	_, err := s.store.User.GetByID(ctx, userID)
	if err != nil {
		switch {
		case err == store.ErrNotFound:
			return apperrors.ErrUserNotFound
		default:
			return fmt.Errorf("failed to verify target user: %w", err)
		}
	}

	if err := s.store.User.UnfollowUser(ctx, userID, followerID); err != nil {
		return err
	}
	return nil
}
