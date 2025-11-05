package usecase

import (
	"context"

	"github.com/team-nino/iam_service/internal/domain/entity"
	"github.com/team-nino/iam_service/internal/domain/repository"
)

// UserUseCase handles user management business logic
type UserUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase creates a new user management use case
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// GetUser retrieves a user by ID
func (u *UserUseCase) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

// ListUsers retrieves a list of users with pagination
func (u *UserUseCase) ListUsers(ctx context.Context, offset, limit int) ([]*entity.User, error) {
	return u.userRepo.List(ctx, offset, limit)
}

// UpdateUser updates user information
func (u *UserUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	return u.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user
func (u *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
}
