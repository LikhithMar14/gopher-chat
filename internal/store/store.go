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
}

type PostRepository interface {
	Create(context.Context, *Post) error
	GetByID(context.Context, int64) (*Post, error)
}

type UserRepository interface {
	Create(context.Context, *User) error
	GetAll(context.Context) ([]User, error)
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Post: &PostStorage{db},
		User: &UserStorage{db},
	}
}
