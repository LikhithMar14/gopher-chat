package handlers

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/LikhithMar14/gopher-chat/internal/errors"
	"errors"	
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
)

type PostHandler struct {
	postService *service.PostService
	commentService *service.CommentService
}

func NewPostHandler(postService *service.PostService, commentService *service.CommentService) *PostHandler {
	return &PostHandler{
		postService: postService,
		commentService: commentService,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req service.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := service.Validate.Struct(req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())	
		return
	}
	ctx := r.Context()
	ctx = utils.SetUserID(ctx, int64(1))
	post, err := h.postService.CreatePost(ctx, req)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"post": post})
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	post, ok := h.postService.GetPostFromContext(ctx)
	if !ok {
		utils.WriteJSONError(w, http.StatusNotFound, "Post not found")
		return
	}
	comments , err := h.commentService.GetCommentsByPostID(ctx, post.ID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	post.Comments = comments
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"post": post})
}
func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ReadIDParam(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	ctx = utils.SetUserID(ctx, int64(1))
	if err := h.postService.DeletePost(ctx, id); err != nil {
		
		switch {
			case err == apperrors.ErrPostNotFound:
				utils.WriteJSONError(w, http.StatusNotFound, apperrors.ErrPostNotFound.Error())
				return
			default:
				utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Post deleted successfully"})
}

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request){

	var req service.UpdatePostRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := service.Validate.Struct(req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := h.postService.UpdatePost(r.Context(), req)
	if err != nil {
		switch {
			case errors.Is(err, apperrors.ErrPostNotFound):
			utils.WriteJSONError(w, http.StatusNotFound, apperrors.ErrPostNotFound.Error())
			return
		default:
			utils.WriteJSONError(w, http.StatusInternalServerError, apperrors.ErrPostNotFound.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"post": post})
	
}

