package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
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
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to get users")
		return
	}

	data := utils.Envelope{
		"users": users,
		"count": len(users),
	}

	if err := utils.WriteJSON(w, http.StatusOK, data); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to write response")
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req service.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	user, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	data := utils.Envelope{
		"user":    user,
		"message": "User created successfully",
	}

	if err := utils.WriteJSON(w, http.StatusCreated, data); err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to write response")
	}
}
