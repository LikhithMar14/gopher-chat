package service

import (
	"context"
	"log"

	"github.com/LikhithMar14/gopher-chat/internal/errors"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	"github.com/LikhithMar14/gopher-chat/internal/utils"

)

type CreatePostRequest struct {
	Title   string   `json:"title" validate:"required,min=3,max=100"`
	Content string   `json:"content" validate:"required,min=10,max=1000"`
	Tags    []string `json:"tags" validate:"required,min=1,max=5"`
}
type UpdatePostRequest struct {
	Title   *string   `json:"title" validate:"omitempty,max=100"`
	Content *string   `json:"content" validate:"omitempty,max=1000"`
	Tags    *[]string `json:"tags" validate:"omitempty,max=5"`
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
	if err := Validate.Struct(req); err != nil {
		log.Print("Validation error",err)
		return nil, err
	}
	if !ok {
		log.Println("user_id not found in context")
		return nil, errors.ErrUserIDNotFound
	}
	log.Println("userID",userID)
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
		switch{
		case err == errors.ErrNotFound:
			log.Println("Error from Service GetPostByID",err)
			return nil, errors.ErrPostNotFound
		default:
			return nil, err
		}
	}
	return post, nil
}
func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	if err := s.store.Post.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}


func (s *PostService) UpdatePost(ctx context.Context, req UpdatePostRequest) (*store.Post, error) {
	post, ok := s.GetPostFromContext(ctx)

	if !ok {
		return nil, errors.ErrPostNotFound
	}
	if req.Title != nil {
	post.Title = *req.Title
	}
	if req.Content != nil {
		post.Content = *req.Content
	}
	if req.Tags != nil {
		post.Tags = *req.Tags
	}
	if err := s.store.Post.Update(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (s *PostService) GetPostFromContext(ctx context.Context) (*store.Post, bool) {
	post, ok := ctx.Value(utils.PostIDKey).(*store.Post)
	return post, ok
}
