package store

import (
	"context"
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64
	PostID    int64
	UserID    int64
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	User 	 User
}

type CommentStorage struct {
	db *sql.DB
}

func (s *CommentStorage) Create(ctx context.Context, comment *Comment) (*Comment, error) {
	return nil, nil
}

func (s *CommentStorage) GetByPostID(ctx context.Context, postID int64) ([]*Comment, error) {
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

	comments := []*Comment{}

	for rows.Next(){
		var c Comment
		c.User = User{}

		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	return comments, nil
}
