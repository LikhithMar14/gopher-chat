package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/lib/pq"
)


type PostStorage struct {
	db *sql.DB
}

func (s *PostStorage) Create(ctx context.Context, post *models.Post) error {
	query := `INSERT INTO posts (content, title, user_id ,tags)
	 VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at, version`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	if err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, post.UserID, post.Tags).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.Version); err != nil {
		return err
	}

	return nil

}

func (s *PostStorage) GetByID(ctx context.Context, id int64) (*models.Post, error) {
	var post models.Post
	query := `
		SELECT id, user_id, title, content, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1

		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	if err := s.db.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt, &post.Version); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStorage) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return apperrors.ErrPostNotFound
	}
	return nil
}

func (s *PostStorage) Update(ctx context.Context, post *models.Post) error {
	query := `
	UPDATE posts
	SET title=$1, content=$2, tags=$3, updated_at=$4, version=version+1
	WHERE id=$5 AND version=$6
	RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	err := s.db.QueryRowContext(ctx, query, post.Title, post.Content, pq.Array(post.Tags), time.Now(), post.ID, post.Version).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			// Check if post exists at all
			var exists bool
			checkQuery := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`
			if checkErr := s.db.QueryRowContext(ctx, checkQuery, post.ID).Scan(&exists); checkErr != nil {
				return checkErr
			}
			if !exists {
				return apperrors.ErrPostNotFound
			}
			// Post exists but version doesn't match - version conflict
			return apperrors.ErrVersionConflict
		default:
			return err
		}
	}

	return nil
}

// UpdateWithOptimisticLocking fetches the latest version and applies updates
func (s *PostStorage) UpdateWithOptimisticLocking(ctx context.Context, id int64, updateFn func(*models.Post) error) (*models.Post, error) {
	// Start a transaction for consistency
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch the latest version with FOR UPDATE to prevent concurrent modifications
	var post models.Post
	query := `
		SELECT id, user_id, title, content, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
		FOR UPDATE
	`

	if err := tx.QueryRowContext(ctx, query, id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt, &post.Version); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, apperrors.ErrPostNotFound
		default:
			return nil, err
		}
	}

	// Apply the updates
	if err := updateFn(&post); err != nil {
		return nil, err
	}

	// Update with the fresh version
	updateQuery := `
		UPDATE posts
		SET title=$1, content=$2, tags=$3, updated_at=$4, version=version+1
		WHERE id=$5 AND version=$6
		RETURNING version, updated_at
	`

	if err := tx.QueryRowContext(ctx, updateQuery, post.Title, post.Content, pq.Array(post.Tags), time.Now(), post.ID, post.Version).Scan(&post.Version, &post.UpdatedAt); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, apperrors.ErrVersionConflict
		default:
			return nil, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &post, nil
}
