package store

import (
	"context"
	"database/sql"

	"github.com/LikhithMar14/gopher-chat/internal/models"
)

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	if err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, []byte(user.PasswordHash)).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetAll(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
