package handlers

import (
	"encoding/json"
	"net/http"

	"errors"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
	"go.uber.org/zap"
)

type PostHandler struct {
	postService    *service.PostService
	commentService *service.CommentService
	logger         *zap.SugaredLogger
}

func NewPostHandler(postService *service.PostService, commentService *service.CommentService, logger *zap.SugaredLogger) *PostHandler {
	return &PostHandler{
		postService:    postService,
		commentService: commentService,
		logger:         logger,
	}
}

// CreatePost godocn m   
//
//	@Summary		Create a new post
//	@Description	Create a new post with title, content, and tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			post	body		models.CreatePostRequest	true	"Post creation request"
//	@Success		201		{object}	map[string]interface{}		"Post created successfully"
//	@Failure		400		{object}	map[string]interface{}		"Validation error"
//	@Failure		500		{object}	map[string]interface{}		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts [post]
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Context-based logging approach (alternative):
	// logger := utils.GetLogger(r.Context())
	// logger.Info("Inside Create Post Handler")

	// Current approach using struct field:
	h.logger.Info("Inside Create Post Handler")

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}
	if err := service.Validate.Struct(req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	ctx := r.Context()
	ctx = utils.SetUserID(ctx, int64(688))
	// will return internal server error if user id is not valid

	post, err := h.postService.CreatePost(ctx, req)
	if err != nil {
		utils.HandleInternalError(w, err)
		return
	}

	data := map[string]interface{}{
		"post": post,
	}
	utils.WriteSuccessResponse(w, http.StatusCreated, data)
}

// GetPostByID godoc
//
//	@Summary		Get post by ID
//	@Description	Retrieve a specific post by its ID, including comments
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"Post ID"
//	@Success		200	{object}	map[string]interface{}	"Post retrieved successfully"
//	@Failure		404	{object}	map[string]interface{}	"Post not found"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Router			/posts/{id} [get]
func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post, ok := h.postService.GetPostFromContext(ctx)

	if !ok {
		utils.WriteErrorResponse(w, http.StatusNotFound, "Post not found")
		return
	}
	comments, err := h.commentService.GetCommentsByPostID(ctx, post.ID)
	if err != nil {
		utils.HandleInternalError(w, err)
		return
	}
	post.Comments = comments

	data := map[string]interface{}{
		"post": post,
	}
	utils.WriteSuccessResponse(w, http.StatusOK, data)
}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete a post by its ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"Post ID"
//	@Success		200	{object}	map[string]interface{}	"Post deleted successfully"
//	@Failure		404	{object}	map[string]interface{}	"Post not found"
//	@Failure		500	{object}	map[string]interface{}	"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [delete]
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	// Set user ID in context (hardcoded for demo - would come from auth middleware in real app)
	ctx := r.Context()
	ctx = utils.SetUserID(ctx, int64(1))
	if err := h.postService.DeletePost(ctx, id); err != nil {
		switch {
		case errors.Is(err, apperrors.ErrPostNotFound):
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		default:
			utils.HandleInternalError(w, err)
		}
		return
	}

	data := map[string]interface{}{
		"message": "Post deleted successfully",
	}
	utils.WriteSuccessResponse(w, http.StatusOK, data)
}

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Update an existing post's title, content, or tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Post ID"
//	@Param			post	body		models.UpdatePostRequest	true	"Post update request"
//	@Success		200		{object}	map[string]interface{}		"Post updated successfully"
//	@Failure		400		{object}	map[string]interface{}		"Validation error"
//	@Failure		404		{object}	map[string]interface{}		"Post not found"
//	@Failure		409		{object}	map[string]interface{}		"Version conflict"
//	@Failure		500		{object}	map[string]interface{}		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{id} [patch]
func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var req models.UpdatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}
	if err := service.Validate.Struct(req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	post, err := h.postService.UpdatePost(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrPostNotFound):
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		case errors.Is(err, apperrors.ErrVersionConflict):
			utils.WriteErrorResponse(w, http.StatusConflict, err.Error())
		default:
			utils.HandleInternalError(w, err)
		}
		return
	}

	data := map[string]interface{}{
		"post": post,
	}
	utils.WriteSuccessResponse(w, http.StatusOK, data)
}

func (h *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}
	if err := service.Validate.Struct(req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	// Set user ID in context (hardcoded for demo - would come from auth middleware in real app)
	ctx := r.Context()
	ctx = utils.SetUserID(ctx, int64(1))

	comment, err := h.commentService.CreateComment(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrPostNotFound):
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		case errors.Is(err, apperrors.ErrUserIDNotFound):
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		default:
			utils.HandleInternalError(w, err)
		}
		return
	}

	data := map[string]interface{}{
		"comment": comment,
	}
	utils.WriteSuccessResponse(w, http.StatusCreated, data)
}
