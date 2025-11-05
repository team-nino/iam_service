package repository

import (
	"context"

	"github.com/team-nino/iam_service/internal/domain/entity"
)

// SessionRepository defines the interface for session data access
type SessionRepository interface {
	// Create creates a new session
	Create(ctx context.Context, session *entity.Session) error

	// GetByToken retrieves a session by its token
	GetByToken(ctx context.Context, token string) (*entity.Session, error)

	// GetByUserID retrieves all sessions for a user
	GetByUserID(ctx context.Context, userID string) ([]*entity.Session, error)

	// Delete deletes a session by its token
	Delete(ctx context.Context, token string) error

	// DeleteByUserID deletes all sessions for a user
	DeleteByUserID(ctx context.Context, userID string) error

	// DeleteExpired deletes all expired sessions
	DeleteExpired(ctx context.Context) error
}
