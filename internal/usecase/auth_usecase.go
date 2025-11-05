package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
	"github.com/team-nino/iam_service/internal/domain/entity"
	"github.com/team-nino/iam_service/internal/domain/repository"
)

// AuthUseCase handles authentication business logic
// This follows the Single Responsibility Principle
type AuthUseCase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

// NewAuthUseCase creates a new authentication use case
func NewAuthUseCase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// Register registers a new user
func (a *AuthUseCase) Register(ctx context.Context, email, username, password string) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := a.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, entity.ErrUserExists
	}

	existingUser, err = a.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, entity.ErrUserExists
	}

	// Create new user
	user, err := entity.NewUser(email, username, password)
	if err != nil {
		return nil, err
	}

	// Generate UUID for the user
	user.ID = uuid.New().String()

	// Save user
	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and creates a session
func (a *AuthUseCase) Login(ctx context.Context, emailOrUsername, password string) (*entity.Session, error) {
	// Try to find user by email first
	user, err := a.userRepo.GetByEmail(ctx, emailOrUsername)
	if err != nil || user == nil {
		// If not found by email, try username
		user, err = a.userRepo.GetByUsername(ctx, emailOrUsername)
		if err != nil || user == nil {
			return nil, entity.ErrUserNotFound
		}
	}

	// Check password
	if !user.CheckPassword(password) {
		return nil, entity.ErrInvalidPassword
	}

	// Generate session token
	token, err := generateToken()
	if err != nil {
		return nil, err
	}

	// Create session (expires in 24 hours)
	session := entity.NewSession(user.ID, token, time.Now().Add(24*time.Hour))
	session.ID = uuid.New().String()

	if err := a.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// Logout invalidates a user session
func (a *AuthUseCase) Logout(ctx context.Context, token string) error {
	return a.sessionRepo.Delete(ctx, token)
}

// ValidateSession validates a session token
func (a *AuthUseCase) ValidateSession(ctx context.Context, token string) (*entity.Session, error) {
	session, err := a.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		// Clean up expired session
		_ = a.sessionRepo.Delete(ctx, token)
		return nil, entity.ErrSessionExpired
	}

	return session, nil
}

// generateToken generates a random token for sessions
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
