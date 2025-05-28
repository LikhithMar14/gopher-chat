package service

import (
	"context"
	"math"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
)

type FeedService struct {
	store store.Storage
}

func NewFeedService(store store.Storage) *FeedService {
	return &FeedService{
		store: store,
	}
}

func (s *FeedService) GetUserFeed(ctx context.Context, req models.FeedRequest) (*models.FeedResponse, error) {
	// Get user ID from context
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		return nil, apperrors.ErrUserIDNotFound
	}

	// Validate request
	if err := Validate.Struct(req); err != nil {
		return nil, err
	}

	// Set default values if not provided
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 50 {
		req.PageSize = 50
	}

	// Get feed items from store
	feedItems, totalCount, err := s.store.Post.GetFeed(ctx, userID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := int(math.Ceil(float64(totalCount) / float64(req.PageSize)))

	paginationInfo := &models.PaginationInfo{
		CurrentPage:  req.Page,
		ItemsPerPage: req.PageSize,
		TotalItems:   totalCount,
		TotalPages:   totalPages,
		HasNext:      req.Page < totalPages,
		HasPrevious:  req.Page > 1,
	}

	return &models.FeedResponse{
		Items:      feedItems,
		Pagination: paginationInfo,
	}, nil
}
