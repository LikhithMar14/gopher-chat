package store

import (
	"context"
	"database/sql"
	"time"

	apperrors "github.com/LikhithMar14/gopher-chat/internal/errors"
	"github.com/LikhithMar14/gopher-chat/internal/models"
)

var (
	ErrNotFound = apperrors.ErrNotFound
)
const QueryTimeoutDuration = 5 * time.Second

type Storage struct {
	Post    PostRepository
	User    UserRepository
	Comment CommentRepository
}

type PostRepository interface {
	Create(context.Context, *models.Post) error
	GetByID(context.Context, int64) (*models.Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *models.Post) error
	UpdateWithOptimisticLocking(context.Context, int64, func(*models.Post) error) (*models.Post, error)
}

type UserRepository interface {
	Create(context.Context, *models.User) error
	GetAll(context.Context) ([]models.User, error)
}
type CommentRepository interface {
	Create(context.Context, *models.Comment) (*models.Comment, error)
	GetByPostID(context.Context, int64) ([]*models.Comment, error)
	
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Post:    &PostStorage{db},
		User:    &UserStorage{db},
		Comment: &CommentStorage{db},
	}
}
