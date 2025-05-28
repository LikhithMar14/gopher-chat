package service

import (
	"context"

	"github.com/LikhithMar14/gopher-chat/internal/store"
)

type FollowService struct {
	store store.Storage
}

func NewFollowService(store store.Storage) *FollowService {
	return &FollowService{store: store}
}

func (s *FollowService) FollowUser(ctx context.Context, currentUserID, userID int64) error {
	return s.store.Follow.FollowUser(ctx, currentUserID, userID)
}

func (s *FollowService) UnfollowUser(ctx context.Context, currentUserID, userID int64) error {
	return s.store.Follow.UnfollowUser(ctx, currentUserID, userID)
}

func (s *FollowService) GetFollowerCount(ctx context.Context, userID int64) (int64, error) {
	return s.store.Follow.GetFollowerCount(ctx, userID)
}

func (s *FollowService) GetFollowingCount(ctx context.Context, userID int64) (int64, error) {
	return s.store.Follow.GetFollowingCount(ctx, userID)
}

func (s *FollowService) IsFollowing(ctx context.Context, currentUserID, userID int64) (bool, error) {
	return s.store.Follow.IsFollowing(ctx, currentUserID, userID)
}
