package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
	"github.com/google/uuid"
)

type AuthService struct {
	store          store.Storage
	mailExpiration time.Duration
}

func NewAuthService(store store.Storage, mailExpiration time.Duration) *AuthService {
	return &AuthService{
		store:          store,
		mailExpiration: mailExpiration,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, req models.RegisterUserRequest) (*models.User, error) {
	hashedPassword, err := models.NewPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	plainToken := uuid.New().String()
	fmt.Printf("Plain token (UUID): %s (length: %d)\n", plainToken, len(plainToken))

	hash := sha256.New()
	hash.Write([]byte(plainToken))
	hashedBytes := hash.Sum(nil)
	token := hex.EncodeToString(hashedBytes)

	fmt.Printf("Hashed token: %s (length: %d)\n", token, len(token))
	fmt.Printf("Hash bytes length: %d\n", len(hashedBytes))

	err = s.store.Auth.CreateAndInvite(ctx, user, token, s.mailExpiration)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Create(ctx context.Context, req models.RegisterUserRequest) (*models.User, error) {
	hashedPassword, err := models.NewPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}
	err = s.store.Auth.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {

	// user, err := s.store.User.GetByEmail(ctx, email)
	// if err != nil {
	//     return nil, err
	// }

	// // Check if the password matches
	// if err := user.Password.Compare(password); err != nil {
	//     return nil, fmt.Errorf("invalid credentials")
	// }

	// return user, nil

	return nil, fmt.Errorf("login method needs GetByEmail implementation in UserRepository")
}

func (s *AuthService) ActivateUser(ctx context.Context, token string) error {
	user, err := s.store.Auth.GetUserFromInvitationToken(ctx, token)
	if err != nil {
		return err
	}

	if user.Activated {
		return apperrors.ErrUserAlreadyActivated
	}

	if err := s.store.Auth.ActivateUser(ctx, user.ID); err != nil {
		return err
	}

	if err := s.store.Auth.DeleteInvitationToken(ctx, token); err != nil {
		return err
	}

	return nil
}
