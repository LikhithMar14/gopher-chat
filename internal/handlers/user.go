package handlers

import (
	"encoding/json"
	"errors"
	"log"
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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	user, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrValidation):
			utils.HandleValidationError(w, err)
		default:
			utils.HandleInternalError(w, err)
		}
		return
	}

	data := map[string]interface{}{
		"user":    user,
		"message": "User created successfully",
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, data)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("ctx", ctx)
	user, ok := h.userService.GetUserFromContext(ctx)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}
	data := map[string]interface{}{
		"user": user,
		"message": "User fetched successfully",
		"success": true,
	}

	utils.WriteSuccessResponse(w, http.StatusOK, data)
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	
}

func (h *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	
}