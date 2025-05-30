package store

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
	"github.com/lib/pq"
)

type AuthStorage struct {
	db *sql.DB
}

func (s *AuthStorage) Create(ctx context.Context, user *models.User) error {
	return s.createUser(ctx, s.db, user)
}

func (s *AuthStorage) createUser(ctx context.Context, executor interface{}, user *models.User) error {

	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var err error
	switch e := executor.(type) {

	case *sql.DB:
		err = e.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.Hash()).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	case *sql.Tx:
		err = e.QueryRowContext(ctx, query, user.Username, user.Email, user.Password.Hash()).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	}

	if err != nil {
		log.Println("Error creating user", err)
		return s.handleUserCreationError(err)
	}
	log.Println("User created successfully", user)
	return nil
}

func (s *AuthStorage) handleUserCreationError(err error) error {
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == "23505" {
			constraint := pqErr.Constraint
			switch constraint {
			case "users_username_key":
				return apperrors.ErrUsernameTaken
			case "users_email_key":
				return apperrors.ErrEmailTaken
			}

			if strings.Contains(pqErr.Message, "username") {
				return apperrors.ErrUsernameTaken
			}
			if strings.Contains(pqErr.Message, "email") {
				return apperrors.ErrEmailTaken
			}

			return apperrors.ErrUserAlreadyExists
		}
	}

	return err
}

func (s *AuthStorage) CreateAndInvite(ctx context.Context, user *models.User, token string, invitationExp time.Duration) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		log.Println("Inside CreateAndInvite")
		if err := s.createUser(ctx, tx, user); err != nil {
			return err
		}

		if err := s.createUserInvitation(ctx, tx, token, invitationExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthStorage) createUserInvitation(ctx context.Context, tx *sql.Tx, token string, invitationExp time.Duration, userID int64) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES (decode($1, 'hex'), $2, $3)`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthStorage) GetUserFromInvitationToken(ctx context.Context, hashedToken string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.activated, u.created_at, u.updated_at, ui.expiry
		FROM users u
		INNER JOIN user_invitations ui ON u.id = ui.user_id
		WHERE encode(ui.token, 'hex') = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	log.Printf("Looking up hashed token: %s (length: %d)", hashedToken, len(hashedToken))

	var user models.User
	var expiry time.Time
	err := s.db.QueryRowContext(ctx, query, hashedToken).Scan(
		&user.ID, &user.Username, &user.Email, &user.Activated,
		&user.CreatedAt, &user.UpdatedAt, &expiry)

	if err != nil {
		log.Printf("Error in token lookup: %v", err)
		switch {
		case err == sql.ErrNoRows:
			return nil, apperrors.ErrInvalidToken
		default:
			return nil, err
		}
	}

	if time.Now().After(expiry) {
		log.Printf("Token expired: %v > %v", time.Now(), expiry)
		return nil, apperrors.ErrTokenExpired
	}

	log.Printf("Found user: %d, activated: %v", user.ID, user.Activated)
	return &user, nil
}

func (s *AuthStorage) ActivateUser(ctx context.Context, userID int64) error {
	query := `UPDATE users SET activated = true WHERE id = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID)
	return err
}

func (s *AuthStorage) DeleteInvitationToken(ctx context.Context, hashedToken string) error {
	query := `DELETE FROM user_invitations WHERE encode(token, 'hex') = $1`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, hashedToken)
	return err
}
