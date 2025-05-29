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
//	@Description	Follow another user to see their posts in your feed
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"User ID to follow"
//	@Success		200	{object}	utils.StandardResponse	"User followed successfully"
//	@Failure		404	{object}	utils.StandardResponse	"User not found"
//	@Failure		500	{object}	utils.StandardResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/follow [post]
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

	data := map[string]interface{}{
		"message": "User followed successfully",
	}
	utils.WriteSuccessResponse(w, http.StatusOK, data)
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user to stop seeing their posts in your feed
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"User ID to unfollow"
//	@Success		200	{object}	utils.StandardResponse	"User unfollowed successfully"
//	@Failure		404	{object}	utils.StandardResponse	"User not found"
//	@Failure		500	{object}	utils.StandardResponse	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/users/{id}/follow [delete]
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

	data := map[string]interface{}{
		"message": "User unfollowed successfully",
	}
	utils.WriteSuccessResponse(w, http.StatusOK, data)
}
