package store

import (
	"context"
	"database/sql"
)

type FollowStorage struct {
	db *sql.DB
}

func (s *FollowStorage) FollowUser(ctx context.Context, currentUserID int64, userID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2) ON CONFLICT (user_id, follower_id) DO NOTHING`
	if _, err := s.db.ExecContext(ctx, query, currentUserID, userID); err != nil {
		return err
	}
	return nil
}

func (s *FollowStorage) UnfollowUser(ctx context.Context, currentUserID int64, userID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`
	if _, err := s.db.ExecContext(ctx, query, currentUserID, userID); err != nil {
		return err
	}
	return nil
}

func (s *FollowStorage) GetFollowerCount(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM followers WHERE user_id = $1`
	var count int64
	if err := s.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *FollowStorage) GetFollowingCount(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM followers WHERE follower_id = $1`
	var count int64
	if err := s.db.QueryRowContext(ctx, query, userID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *FollowStorage) IsFollowing(ctx context.Context, currentUserID, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM followers WHERE user_id = $1 AND follower_id = $2)`
	var exists bool
	if err := s.db.QueryRowContext(ctx, query, userID, currentUserID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
