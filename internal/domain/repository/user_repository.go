package repository

import (
	"context"

	"github.com/team-nino/iam_service/internal/domain/entity"
)

// UserRepository defines the interface for user data access
// This follows the Dependency Inversion Principle - high-level modules
// don't depend on low-level modules, both depend on abstractions
type UserRepository interface {
	// Create creates a new user in the repository
	Create(ctx context.Context, user *entity.User) error

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*entity.User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetByUsername retrieves a user by their username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// Delete deletes a user by their ID
	Delete(ctx context.Context, id string) error

	// List retrieves all users with pagination
	List(ctx context.Context, offset, limit int) ([]*entity.User, error)
}
