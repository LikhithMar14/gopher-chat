package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
)

var (
	ErrNotFound = apperrors.ErrNotFound
)

const QueryTimeoutDuration = 3 * time.Second

type Storage struct {
	Post    PostRepository
	User    UserRepository
	Comment CommentRepository
	Follow  FollowRepository
	Auth    AuthRepository
}

type PostRepository interface {
	Create(context.Context, *models.Post) error
	GetByID(context.Context, int64) (*models.Post, error)
	Delete(context.Context, int64) error
	Update(context.Context, *models.Post) error
	UpdateWithOptimisticLocking(context.Context, int64, func(*models.Post) error) (*models.Post, error)
	GetFeed(context.Context, int64, int, int) ([]*models.FeedItem, int64, error)
}

type UserRepository interface {
	Create(context.Context, *models.User) error
	GetAll(context.Context) ([]models.User, error)
	GetByID(context.Context, int64) (*models.User, error)
	FollowUser(context.Context, int64, int64) error
	UnfollowUser(context.Context, int64, int64) error
}
type CommentRepository interface {
	Create(context.Context, *models.Comment) (*models.Comment, error)
	GetByPostID(context.Context, int64) ([]*models.Comment, error)
}

type FollowRepository interface {
	FollowUser(context.Context, int64, int64) error
	UnfollowUser(context.Context, int64, int64) error
	GetFollowerCount(context.Context, int64) (int64, error)
	GetFollowingCount(context.Context, int64) (int64, error)
	IsFollowing(context.Context, int64, int64) (bool, error)
}

type AuthRepository interface {
	Create(context.Context, *models.User) error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Post:    &PostStorage{db},
		User:    &UserStorage{db},
		Comment: &CommentStorage{db},
		Follow:  &FollowStorage{db},
	}
}
