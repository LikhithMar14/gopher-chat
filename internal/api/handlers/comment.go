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

type CommentHandler struct {
	commentService *service.CommentService
	postService    *service.PostService
}

func NewCommentHandler(commentService *service.CommentService, postService *service.PostService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		postService:    postService,
	}
}

// CreateComment godoc
//
//	@Summary		Create a comment on a post
//	@Description	Add a new comment to a specific post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int							true	"Post ID"
//	@Param			comment	body		models.CreateCommentRequest	true	"Comment creation request"
//	@Success		201		{object}	utils.StandardResponse		"Comment created successfully"
//	@Failure		400		{object}	utils.StandardResponse		"Validation error"
//	@Failure		401		{object}	utils.StandardResponse		"Unauthorized"
//	@Failure		404		{object}	utils.StandardResponse		"Post not found"
//	@Failure		500		{object}	utils.StandardResponse		"Internal server error"
//	@Security		ApiKeyAuth
//	@Router			/posts/{id}/comments [post]
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	if err := service.Validate.Struct(req); err != nil {
		utils.HandleValidationError(w, err)
		return
	}

	ctx := r.Context()
	ctx = utils.SetUserID(ctx, int64(668))
	//will give internal server error if the userid is not correct

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

// GetCommentsByPostID godoc
//
//	@Summary		Get comments for a post
//	@Description	Retrieve all comments for a specific post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int						true	"Post ID"
//	@Success		200	{object}	utils.StandardResponse	"Comments retrieved successfully"
//	@Failure		404	{object}	utils.StandardResponse	"Post not found"
//	@Failure		500	{object}	utils.StandardResponse	"Internal server error"
//	@Router			/posts/{id}/comments [get]
func (h *CommentHandler) GetCommentsByPostID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get post from context (set by middleware)
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

	data := map[string]interface{}{
		"comments": comments,
		"count":    len(comments),
		"post_id":  post.ID,
	}
	utils.WriteSuccessResponse(w, http.StatusOK, data)
}
