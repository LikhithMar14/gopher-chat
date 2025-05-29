package handlers

import (
	"net/http"
	"strconv"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
)

type FeedHandler struct {
	userService *service.UserService
	postService *service.PostService
	feedService *service.FeedService
}

func NewFeedHandler(userService *service.UserService, postService *service.PostService, feedService *service.FeedService) *FeedHandler {
	return &FeedHandler{
		userService: userService,
		postService: postService,
		feedService: feedService,
	}
}

// GetFeed godoc
//
//	@Summary		Get user feed
//	@Description	Retrieve a paginated feed of posts from followed users
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int						false	"Page number (default: 1)"
//	@Param			page_size	query		int						false	"Items per page (default: 10, max: 50)"
//	@Success		200			{object}	models.FeedResponse		"Feed retrieved successfully"
//	@Failure		500			{object}	utils.StandardResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/users/me/feed [get]
func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// if _, ok := h.userService.GetUserFromContext(ctx); !ok {
	// 	utils.WriteErrorResponse(w, http.StatusUnauthorized, "User not found")
	// 	return
	// }

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("page_size")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 50 {
			pageSize = ps
		}
	}

	feedRequest := models.FeedRequest{
		Page:     page,
		PageSize: pageSize,
	}

	feedResponse, err := h.feedService.GetUserFeed(ctx, feedRequest)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve feed")
		return
	}

	utils.WriteSuccessResponse(w, http.StatusOK, feedResponse)
}
