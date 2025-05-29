package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	text *string
	hash []byte
}

func NewPassword(plaintext string) (*Password, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Password{
		text: &plaintext,
		hash: hash,
	}, nil
}

func NewPasswordFromHash(hash []byte) *Password {
	return &Password{
		text: nil,
		hash: hash,
	}
}

func (p *Password) Compare(plaintext string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
}

func (p *Password) Hash() []byte {
	return p.hash
}

func (p *Password) Text() *string {
	return p.text
}

func (p *Password) String() string {
	if p.text != nil {
		return *p.text
	}
	return ""
}

func (p *Password) HasText() bool {
	return p.text != nil
}

type Post struct {
	ID        int64      `json:"id"`
	Content   string     `json:"content"`
	Title     string     `json:"title"`
	UserID    int64      `json:"user_id"`
	Tags      []string   `json:"tags"`
	Version   int        `json:"version"`
	Comments  []*Comment `json:"comments"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user"`
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  *Password `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

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

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}
type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1,max=500"`
}

type FollowUnfollowRequest struct {
	UserID int64 `json:"user_id"`
}

type FeedItem struct {
	Post   *Post `json:"post"`
	Author *User `json:"author"`
}

type FeedResponse struct {
	Items      []*FeedItem     `json:"items"`
	Pagination *PaginationInfo `json:"pagination"`
}

type PaginationInfo struct {
	CurrentPage  int   `json:"current_page"`
	ItemsPerPage int   `json:"items_per_page"`
	TotalItems   int64 `json:"total_items"`
	TotalPages   int   `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

type FeedRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=50"`
}
