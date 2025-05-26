package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("record not found")
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
