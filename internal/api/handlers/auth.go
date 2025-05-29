package handlers

import (
	"net/http"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// registerUserHandler godoc
//
//	@Summary		Registers a user
//	@Description	Registers a new user with the provided information
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.RegisterUserRequest	true	"User registration details"
//	@Success		201		{object}	utils.StandardResponse		"User registered successfully"
//	@Failure		400		{object}	utils.StandardResponse		"Invalid request or validation error"
//	@Failure		500		{object}	utils.StandardResponse		"Internal server error"
//	@Router			/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterUserRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	if err := service.Validate.Struct(req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	user, err := h.authService.Register(r.Context(), req)
	if err != nil {
		utils.HandleInternalError(w, err)
		return
	}

	data := map[string]interface{}{
		"user":    user,
		"message": "User registered successfully",
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, data)
}
