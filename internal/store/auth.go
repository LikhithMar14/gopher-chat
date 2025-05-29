package store

import (
	"context"
	"database/sql"
	"github.com/LikhithMar14/gopher-chat/internal/models"
)

type AuthStorage struct {
	db *sql.DB
}


func (s *AuthStorage) Create(ctx context.Context, user *models.User) error {
	return nil
}

