package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

var (
	ErrNotFound     = errors.New("record not found")
	ErrInternal     = errors.New("internal server error")
	ErrInvalidInput = errors.New("invalid input format")
)

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req service.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, ErrInvalidInput.Error())
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
	id, err := utils.ReadIDParam(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid Post ID")
		return
	}
	ctx := r.Context()
	ctx = utils.SetUserID(ctx, int64(1))
	post, err := h.postService.GetPostByID(ctx, id)
	if err != nil {
		switch {
		case err == ErrNotFound:
			utils.WriteJSONError(w, http.StatusNotFound, "Post not found")
		default:
			utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"post": post})
}
