package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/LikhithMar14/gopher-chat/internal/config"
	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/store"
	apperrors "github.com/LikhithMar14/gopher-chat/internal/utils/errors"
	"github.com/LikhithMar14/gopher-chat/internal/utils/mailer"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService struct {
	store          store.Storage
	mailExpiration time.Duration
	mailer mailer.Client
	frontendURL string
	logger *zap.SugaredLogger
	config config.Config
	}

func NewAuthService(store store.Storage, mailExpiration time.Duration, mailer mailer.Client, config config.Config, logger *zap.SugaredLogger) *AuthService {

	return &AuthService{
		store:          store,
		mailExpiration: mailExpiration,
		mailer:         mailer,
		frontendURL:    config.FrontendURL,
		logger:         logger,
		config:         config,
	}
}

// generatePlainToken generates a new UUID-based token
func (s *AuthService) generatePlainToken() string {
	return uuid.New().String()
}

// hashToken hashes a plain token using SHA256
func (s *AuthService) hashToken(plainToken string) string {
	hash := sha256.New()
	hash.Write([]byte(plainToken))
	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

func (s *AuthService) RegisterUser(ctx context.Context, req models.RegisterUserRequest) (*models.User, string, error) {
	hashedPassword, err := models.NewPassword(req.Password)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Generate plain token to give to user
	plainToken := s.generatePlainToken()
	fmt.Printf("Plain token (UUID): %s (length: %d)\n", plainToken, len(plainToken))

	hashedToken := s.hashToken(plainToken)

	fmt.Printf("Hashed token for storage: %s (length: %d)\n", hashedToken, len(hashedToken))

	err = s.store.Auth.CreateAndInvite(ctx, user, hashedToken, s.mailExpiration)
	if err != nil {
		return nil, "", err
	}
	// sending email to user
	isProdEnv := s.config.Env == "prod"
	s.logger.Info("Sending email to user", zap.String("email", user.Email), zap.String("username", user.Username), zap.String("frontendURL", s.frontendURL), zap.String("hashedToken", hashedToken), zap.Bool("isProdEnv", isProdEnv))

	vars := map[string]interface{}{
		"Username": user.Username,
		"ActivationURL": fmt.Sprintf("%s/activate/%s", s.frontendURL, hashedToken),
	}

	if err := s.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv); err != nil {
		// SAGA PATTERN
		if err := s.store.Auth.DeleteInvitationToken(ctx, hashedToken); err != nil {
			s.logger.Error("Error deleting invitation token", zap.Error(err))
		}
		s.logger.Error("Error sending email", zap.Error(err))
		return nil, "", err
	}

	return user, plainToken, nil
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
	// Hash the incoming plain token to compare with stored hash
	hashedToken := s.hashToken(token)

	user, err := s.store.Auth.GetUserFromInvitationToken(ctx, hashedToken)
	if err != nil {
		return err
	}

	if user.Activated {
		return apperrors.ErrUserAlreadyActivated
	}

	if err := s.store.Auth.ActivateUser(ctx, user.ID); err != nil {
		return err
	}

	if err := s.store.Auth.DeleteInvitationToken(ctx, hashedToken); err != nil {
		return err
	}

	return nil
}
