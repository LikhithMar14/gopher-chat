package handlers

import (
	"encoding/json"

	"net/http"

	"errors"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
)

type PostHandler struct {
	postService    *service.PostService
	commentService *service.CommentService
}

func NewPostHandler(postService *service.PostService, commentService *service.CommentService) *PostHandler {
	return &PostHandler{
		postService:    postService,
		commentService: commentService,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePostRequest
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
