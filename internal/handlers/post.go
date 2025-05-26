package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/LikhithMar14/gopher-chat/internal/errors"
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
	id, err := utils.ReadIDParam(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()
	ctx = utils.SetUserID(ctx, int64(1))
	post, err := h.postService.GetPostByID(ctx, id)

	if err != nil {
		switch {
			case err == errors.ErrNotFound:
				utils.WriteJSONError(w, http.StatusNotFound, err.Error())
				return
			default:
				utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"post": post})
}
