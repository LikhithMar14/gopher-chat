package store

import (
	"context"
	"database/sql"

	apperrors "github.com/LikhithMar14/gopher-chat/internal/errors"
)

var (
	ErrNotFound = apperrors.ErrNotFound
)



type Storage struct {
	Post PostRepository
	User UserRepository
	Comment CommentRepository
}

type PostRepository interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *Post) error
	
}

type UserRepository interface {
	Create(context.Context, *User) error
	GetAll(context.Context) ([]User, error)
}
type CommentRepository interface {
	Create(context.Context, *Comment) (*Comment, error)
	GetByPostID(context.Context, int64) ([]*Comment, error)
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Post: &PostStorage{db},
		User: &UserStorage{db},
		Comment: &CommentStorage{db},
	}
}
