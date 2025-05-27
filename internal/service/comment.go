package service

import (
	"context"

	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/store"
)


type CommentService struct {
	store store.Storage
}

func NewCommentService(store store.Storage) *CommentService {
	return &CommentService{store: store}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	return s.store.Comment.Create(ctx,comment)
}

func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID int64) ([]*models.Comment, error) {
	if postID <= 0 {
		return nil, apperrors.NewBadRequestError("PostId should be Valid")
	}
	return s.store.Comment.GetByPostID(ctx,postID)
}