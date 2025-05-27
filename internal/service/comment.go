package service

import (
	"context"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
)

type CommentService struct {
	store store.Storage
}

func NewCommentService(store store.Storage) *CommentService {
	return &CommentService{store: store}
}

func (s *CommentService) CreateComment(ctx context.Context, req *models.CreateCommentRequest) (*models.Comment, error) {
	// Validate the request
	if err := Validate.Struct(req); err != nil {
		return nil, err
	}

	// Get user ID from context
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return nil, apperrors.ErrUserIDNotFound
	}

	// Get post from context (set by middleware)
	post, ok := ctx.Value(utils.PostIDKey).(*models.Post)
	if !ok {
		return nil, apperrors.ErrPostNotFound
	}

	// Create comment
	comment := &models.Comment{
		PostID:  post.ID,
		UserID:  userID,
		Content: req.Content,
	}

	createdComment, err := s.store.Comment.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return createdComment, nil
}

func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID int64) ([]*models.Comment, error) {
	if postID <= 0 {
		return nil, apperrors.NewBadRequestError("PostId should be Valid")
	}
	return s.store.Comment.GetByPostID(ctx, postID)
}
