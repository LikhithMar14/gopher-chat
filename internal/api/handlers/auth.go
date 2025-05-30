package handlers

import (
	"errors"

	"net/http"

	"fmt"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
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
//	@Failure		409		{object}	utils.StandardResponse		"Username or email already exists"
//	@Failure		500		{object}	utils.StandardResponse		"Internal server error"
//	@Router			/auth/register [post]
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RegisterUser")
	var req models.RegisterUserRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	if err := service.Validate.Struct(req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	user, plainToken, err := h.authService.RegisterUser(r.Context(), req)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	data := map[string]interface{}{
		"user":    user,
		"token":   plainToken,
		"message": "User registered successfully. Please check your email for activation.",
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, data)
}

/*
Create User Seed is only for seeding the database with users.
*/
func (h *AuthHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterUserRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	if err := service.Validate.Struct(req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	user, err := h.authService.Create(r.Context(), req)
	if err != nil {
		h.handleAuthError(w, err)
		return
	}

	data := map[string]interface{}{
		"user":    user,
		"message": "User created successfully",
	}

	utils.WriteSuccessResponse(w, http.StatusCreated, data)
}

// handleAuthError provides specific error handling for authentication-related errors
func (h *AuthHandler) handleAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, apperrors.ErrUsernameTaken):
		utils.WriteErrorResponse(w, http.StatusConflict, "Username is already taken")
	case errors.Is(err, apperrors.ErrEmailTaken):
		utils.WriteErrorResponse(w, http.StatusConflict, "Email is already registered")
	case errors.Is(err, apperrors.ErrUserAlreadyExists):
		utils.WriteErrorResponse(w, http.StatusConflict, "User already exists")
	default:
		utils.HandleInternalError(w, err)
	}
}

// ActivateUser godoc
//
//	@Summary		Activates a user
//	@Description	Activates a user with the provided token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			token	path		string					true	"Activation token"
//	@Success		200		{object}	utils.StandardResponse	"User activated successfully"
//	@Failure		400		{object}	utils.StandardResponse	"Invalid token format"
//	@Failure		404		{object}	utils.StandardResponse	"Invalid or expired token"
//	@Failure		409		{object}	utils.StandardResponse	"User already activated"
//	@Failure		500		{object}	utils.StandardResponse	"Internal server error"
//	@Router			/users/activate/{token} [put]
func (h *AuthHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	token := utils.ReadStringParam(r, "token")
	fmt.Printf("Received token: '%s' (length: %d)\n", token, len(token))

	if token == "" {
		utils.HandleValidationError(w, errors.New("token is required"))
		return
	}

	err := h.authService.ActivateUser(r.Context(), token)
	if err != nil {
		fmt.Printf("Activation error: %v\n", err)
		switch {
		case errors.Is(err, apperrors.ErrInvalidToken):
			utils.WriteErrorResponse(w, http.StatusNotFound, "Invalid or expired activation token")
		case errors.Is(err, apperrors.ErrTokenExpired):
			utils.WriteErrorResponse(w, http.StatusNotFound, "Activation token has expired")
		case errors.Is(err, apperrors.ErrUserAlreadyActivated):
			utils.WriteErrorResponse(w, http.StatusConflict, "User is already activated")
		case errors.Is(err, apperrors.ErrUserNotFound):
			utils.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		default:
			utils.HandleInternalError(w, err)
		}
		return
	}

	data := map[string]interface{}{
		"message": "User activated successfully",
	}
	utils.WriteSuccessResponse(w, http.StatusOK, data)
}
