package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/pkg/database"
	apperrors "github.com/LikhithMar14/gopher-chat/pkg/errors"
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
	CreateAndInvite(context.Context, *models.User, string, time.Duration) error
	Create(context.Context, *models.User) error
	GetUserFromInvitationToken(context.Context, string) (*models.User, error)
	ActivateUser(context.Context, int64) error
	DeleteInvitationToken(context.Context, string) error
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Post:    &PostStorage{db},
		User:    &UserStorage{db},
		Comment: &CommentStorage{db},
		Follow:  &FollowStorage{db},
		Auth:    &AuthStorage{db},
	}
}

func withTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	return database.WithTx(ctx, db, fn)
}
