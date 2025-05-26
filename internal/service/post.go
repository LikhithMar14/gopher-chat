package service

import (
	"context"
	"errors"
	"log"

	"github.com/LikhithMar14/gopher-chat/internal/store"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
)

type CreatePostRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type PostService struct {
	store store.Storage
}

func NewPostService(store store.Storage) *PostService {
	return &PostService{
		store: store,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req CreatePostRequest) (*store.Post, error) {
	var post store.Post
	userID, ok := utils.GetUserID(ctx)
	if !ok {
		log.Println("user_id not found in context")
		return nil, errors.New("user_id not found in context")
	}
	post.Title = req.Title
	post.Content = req.Content
	post.UserID = userID
	post.Tags = req.Tags
	if err := s.store.Post.Create(ctx, &post); err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *PostService) GetPostByID(ctx context.Context, id int64) (*store.Post, error) {
	post, err := s.store.Post.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}
