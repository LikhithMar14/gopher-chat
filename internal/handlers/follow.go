package handlers

import (
	"errors"
	"log"
	
	"net/http"

	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
)

type FollowHandler struct {
	followService *service.FollowService		
	userService   *service.UserService
}

func NewFollowHandler(followService *service.FollowService, userService *service.UserService) *FollowHandler {
	return &FollowHandler{followService: followService, userService: userService}
}	

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
	log.Println("err in follow handler", err)
	if err != nil {
		utils.HandleInternalError(w, err)
		return
	}
	utils.WriteSuccessResponse(w, http.StatusNoContent, nil)
}
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
	log.Println("err in unfollow handler", err)
	if err != nil {
		utils.HandleInternalError(w, err)
		return
	}
	utils.WriteSuccessResponse(w, http.StatusNoContent, nil)
}