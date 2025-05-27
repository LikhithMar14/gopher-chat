package service

import (
	"context"
	"errors"

	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
)


type PostService struct {
	store store.Storage
}

func NewPostService(store store.Storage) *PostService {
	return &PostService{
		store: store,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req models.CreatePostRequest) (*models.Post, error) {
	var post models.Post
	userID, ok := utils.GetUserID(ctx)
	if err := Validate.Struct(req); err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperrors.ErrUserIDNotFound
	}
	post.Title = req.Title
	post.Content = req.Content
	post.UserID = userID
	post.Tags = req.Tags
	if err := s.store.Post.Create(ctx, &post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	post, err := s.store.Post.GetByID(ctx, id)
	if err != nil {
		switch {
		case err == apperrors.ErrNotFound:
			return nil, apperrors.ErrPostNotFound
		default:
			return nil, err
		}
	}
	return post, nil
}
func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	if err := s.store.Post.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *PostService) UpdatePost(ctx context.Context, req models.UpdatePostRequest) (*models.Post, error) {
	// Get the post ID from context (set by middleware)
	postFromContext, ok := s.GetPostFromContext(ctx)
	if !ok {
		return nil, apperrors.ErrPostNotFound
	}

	// Use optimistic locking with retry logic
	const maxRetries = 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		post, err := s.store.Post.UpdateWithOptimisticLocking(ctx, postFromContext.ID, func(p *models.Post) error {
			// Apply the updates to the fresh post data
			if req.Title != nil {
				p.Title = *req.Title
			}
			if req.Content != nil {
				p.Content = *req.Content
			}
			if req.Tags != nil {
				p.Tags = *req.Tags
			}
			return nil
		})

		if err != nil {
			lastErr = err
			// If it's a version conflict, retry
			if errors.Is(err, apperrors.ErrVersionConflict) {
				continue
			}
			// For other errors, return immediately
			return nil, err
		}

		// Success - return the updated post
		return post, nil
	}

	// If we exhausted all retries, return the last error
	return nil, lastErr
}

func (s *PostService) GetPostFromContext(ctx context.Context) (*models.Post, bool) {
	post, ok := ctx.Value(utils.PostIDKey).(*models.Post)
	return post, ok
}

