package handlers

import (
	"errors"

	"net/http"

	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	"go.uber.org/zap"
)

type FollowHandler struct {
	followService *service.FollowService
	userService   *service.UserService
	logger        *zap.SugaredLogger
}

func NewFollowHandler(followService *service.FollowService, userService *service.UserService, logger *zap.SugaredLogger) *FollowHandler {
	return &FollowHandler{
		followService: followService,
		userService:   userService,
		logger:        logger,
	}
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follow another user by their ID
//	@Tags			follow
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID to follow"
//	@Success		204	"User followed successfully"
//	@Failure		404	{object}	map[string]interface{}	"User not found"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/follow [put]
func (h *FollowHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.userService.GetUserFromContext(ctx)
	if !ok {
		utils.HandleInternalError(w, errors.New("user not found"))
		return
	}

	userID := user.ID
	currentUserID := 667

	err := h.followService.FollowUser(ctx, int64(currentUserID), int64(userID))
	if err != nil {
		h.logger.Error("Error in follow handler", zap.Error(err))
		utils.HandleInternalError(w, err)
		return
	}
	utils.WriteSuccessResponse(w, http.StatusNoContent, nil)
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by their ID
//	@Tags			follow
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"User ID to unfollow"
//	@Success		204	"User unfollowed successfully"
//	@Failure		404	{object}	map[string]interface{}	"User not found"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/unfollow [put]
func (h *FollowHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.userService.GetUserFromContext(ctx)
	if !ok {
		utils.HandleInternalError(w, errors.New("user not found"))
		return
	}
	userID := user.ID
	currentUserID := 667

	err := h.followService.UnfollowUser(ctx, int64(currentUserID), int64(userID))
	if err != nil {
		h.logger.Error("Error in unfollow handler", zap.Error(err))
		utils.HandleInternalError(w, err)
		return
	}
	utils.WriteSuccessResponse(w, http.StatusNoContent, nil)
}
