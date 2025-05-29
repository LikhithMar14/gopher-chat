package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Retrieve a list of all users in the system
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	utils.StandardResponse	"success response with users array and count"
//	@Failure		500	{object}	utils.StandardResponse	"Internal server error"
//	@Router			/users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.userService.GetUsers(ctx)
	if err != nil {
		utils.HandleInternalError(w, err)
		return
	}

	data := map[string]interface{}{
		"users": users,
		"count": len(users),
	}

	utils.WriteSuccessResponse(w, http.StatusOK, data)
}



// GetUserByID godoc
//
//	@Summary		Get user by ID
//	@Description	Get user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	utils.StandardResponse
//	@Failure		400	{object}	utils.StandardResponse
//	@Failure		404	{object}	utils.StandardResponse
//	@Failure		500	{object}	utils.StandardResponse
//	@Failure		401	{object}	utils.StandardResponse
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.userService.GetUserFromContext(ctx)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}
	data := map[string]interface{}{
		"user":    user,
		"message": "User fetched successfully",
		"success": true,
	}

	utils.WriteSuccessResponse(w, http.StatusOK, data)
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.userService.GetUserFromContext(ctx)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	var req models.FollowUnfollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	// Check if trying to follow themselves
	if user.ID == req.UserID {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Cannot follow yourself")
		return
	}

	err := h.userService.FollowUser(ctx, user.ID, req.UserID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserNotFound):
			utils.WriteErrorResponse(w, http.StatusNotFound, "Target user not found")
		case errors.Is(err, apperrors.ErrConflict):
			utils.WriteErrorResponse(w, http.StatusConflict, "Already following this user")
		default:
			utils.HandleInternalError(w, err)
		}
		return
	}

	data := map[string]interface{}{
		"message": "User followed successfully",
		"success": true,
	}

	utils.WriteSuccessResponse(w, http.StatusOK, data)
}

func (h *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := h.userService.GetUserFromContext(ctx)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	var req models.FollowUnfollowRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	if user.ID == req.UserID {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Cannot unfollow yourself")
		return
	}

	err := h.userService.UnfollowUser(ctx, user.ID, req.UserID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserNotFound):
			utils.WriteErrorResponse(w, http.StatusNotFound, "Target user not found")
		case errors.Is(err, apperrors.ErrNotFound):
			utils.WriteErrorResponse(w, http.StatusNotFound, "Not following this user")
		default:
			utils.HandleInternalError(w, err)
		}
		return
	}

	data := map[string]interface{}{
		"message": "User unfollowed successfully",
		"success": true,
	}

	utils.WriteSuccessResponse(w, http.StatusOK, data)
}
