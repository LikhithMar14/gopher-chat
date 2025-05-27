package store

import (
	"context"
	"database/sql"
	"time"
	"github.com/LikhithMar14/gopher-chat/internal/models"
)



type CommentStorage struct {
	db *sql.DB
}
type CreateCommentRequest struct {
	PostID int64
	UserID int64
	Content string
}
func (s *CommentStorage) Create(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, post_id, user_id, content, created_at, updated_at
	`
	ctx,cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	row := s.db.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content)

	var c models.Comment
	err := row.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}


func (s *CommentStorage) GetByPostID(ctx context.Context, postID int64) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, c.updated_at, u.username, u.id
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;

	`
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*models.Comment{}

	for rows.Next(){
		var c models.Comment
		c.User = models.User{}

		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	return comments, nil
}
