package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Post PostRepository
	User UserRepository
}

type PostRepository interface {
	Create(context.Context, *Post) error
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
