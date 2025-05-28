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
