package store

import (
	"context"
	"database/sql"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
)

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	if err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.Hash()).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, email, activated, created_at, updated_at FROM users ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Activated, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserStorage) GetByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, activated, created_at, updated_at FROM users WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user models.User
	var passwordHash []byte
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.Activated, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	user.Password = models.NewPasswordFromHash(passwordHash)

	return &user, nil
}

func (s *UserStorage) FollowUser(ctx context.Context, userID int64, followerID int64) error {
	query := `INSERT INTO followers (user_id, follower_id) VALUES ($1, $2)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		// Check for duplicate key constraint violation (PostgreSQL specific)
		if err.Error() == "pq: duplicate key value violates unique constraint \"followers_pkey\"" ||
			err.Error() == "UNIQUE constraint failed: followers.user_id, followers.follower_id" {
			return apperrors.ErrConflict
		}
		return err
	}
	return nil
}

func (s *UserStorage) UnfollowUser(ctx context.Context, userID int64, followerID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
