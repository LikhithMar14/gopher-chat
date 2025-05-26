package store

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	if err := s.db.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetAll(ctx context.Context) ([]User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
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
