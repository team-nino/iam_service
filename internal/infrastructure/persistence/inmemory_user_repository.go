package persistence

import (
	"context"
	"sync"

	"github.com/team-nino/iam_service/internal/domain/entity"
	"github.com/team-nino/iam_service/internal/domain/repository"
)

// InMemoryUserRepository implements UserRepository using in-memory storage
// This follows the Liskov Substitution Principle - can be replaced with any other implementation
type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*entity.User
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() repository.UserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*entity.User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user already exists
	for _, u := range r.users {
		if u.Email == user.Email || u.Username == user.Username {
			return entity.ErrUserExists
		}
	}

	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, entity.ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, entity.ErrUserNotFound
}

func (r *InMemoryUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, entity.ErrUserNotFound
}

func (r *InMemoryUserRepository) Update(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return entity.ErrUserNotFound
	}

	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return entity.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepository) List(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	// Simple pagination
	start := offset
	if start > len(users) {
		return []*entity.User{}, nil
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}
