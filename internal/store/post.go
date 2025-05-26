package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)
type Post struct {
	ID int64	`json:"id"`
	Content string `json:"content"`
	Title string `json:"title"`
	UserID int64 `json:"user_id"`
	Tags []string `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type PostStorage struct {
	db *sql.DB
}

func (s *PostStorage) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content, title, user_id ,tags) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`


	if err := s.db.QueryRowContext(ctx,query,post.Content,post.Title,post.UserID,post.Tags).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt); err != nil {
		return err
	}

	return nil

}

func (s *PostStorage)GetByID(ctx context.Context,id int64) (*Post, error){
	var post Post
	query := `
		SELECT id, user_id, title, content, tags, created_at, updated_at
		FROM posts
		WHERE id = $1

	`

	if err := s.db.QueryRowContext(ctx,query,id).Scan(&post.ID, &post.UserID, &post.Title, &post.Content, pq.Array(&post.Tags), &post.CreatedAt, &post.UpdatedAt); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}