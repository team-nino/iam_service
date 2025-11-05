package persistence

import (
	"context"
	"sync"
	"time"

	"github.com/team-nino/iam_service/internal/domain/entity"
	"github.com/team-nino/iam_service/internal/domain/repository"
)

// InMemorySessionRepository implements SessionRepository using in-memory storage
type InMemorySessionRepository struct {
	mu       sync.RWMutex
	sessions map[string]*entity.Session
}

// NewInMemorySessionRepository creates a new in-memory session repository
func NewInMemorySessionRepository() repository.SessionRepository {
	return &InMemorySessionRepository{
		sessions: make(map[string]*entity.Session),
	}
}

func (r *InMemorySessionRepository) Create(ctx context.Context, session *entity.Session) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sessions[session.Token] = session
	return nil
}

func (r *InMemorySessionRepository) GetByToken(ctx context.Context, token string) (*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	session, exists := r.sessions[token]
	if !exists {
		return nil, entity.ErrSessionNotFound
	}
	return session, nil
}

func (r *InMemorySessionRepository) GetByUserID(ctx context.Context, userID string) ([]*entity.Session, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sessions := make([]*entity.Session, 0)
	for _, session := range r.sessions {
		if session.UserID == userID {
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

func (r *InMemorySessionRepository) Delete(ctx context.Context, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.sessions, token)
	return nil
}

func (r *InMemorySessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for token, session := range r.sessions {
		if session.UserID == userID {
			delete(r.sessions, token)
		}
	}
	return nil
}

func (r *InMemorySessionRepository) DeleteExpired(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for token, session := range r.sessions {
		if now.After(session.ExpiresAt) {
			delete(r.sessions, token)
		}
	}
	return nil
}
