package services

import (
	"authService/internal/models"
	"authService/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UserService defines the methods for interacting with user data
type UserService[T any] interface {
	CreateUser(user *models.User) (uuid.UUID, error)
	GetUser(id string) (T, error)
	UpdateUser(user *models.User) (T, error)
	DeleteUser(id string) error
	// CheckUserExists(id string) (bool, error)
}

// userService is the struct that implements the UserService interface
type userService[T any] struct {
	userRepo repository.UserRepository
}

// NewUserService initializes a new UserService instance with a generic type T
func NewUserService[T any]() UserService[T] {
	return &userService[T]{
		userRepo: *repository.NewUserRepository(),
	}
}

// CreateUser creates a new user and returns its ID
func (s *userService[T]) CreateUser(user *models.User) (uuid.UUID, error) {
	ctx := context.Background()
	createdUser, err := s.userRepo.Insert(ctx, *user)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser.ID, nil
}

func (s *userService[T]) GetUser(id string) (T, error) {
	var result T
	ctx := context.Background()

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return result, fmt.Errorf("invalid user ID: %w", err)
	}

	// Fetching the user asynchronously, expecting it to be of type models.User
	user, err := s.userRepo.Query(ctx, parsedID)
	if err != nil {
		return result, fmt.Errorf("failed to get user: %w", err)
	}

	// Convert models.User to type T
	result, ok := any(user).(T)
	if !ok {
		return result, fmt.Errorf("failed to convert user to type T")
	}
	return result, nil
}
// UpdateUser updates the user information
func (s *userService[T]) UpdateUser(user *models.User) (T, error) {
	var updatedUser T
	ctx := context.Background()

	// Update user in the repository
	result,err := s.userRepo.Update(ctx, *user)
	if err != nil {
		return updatedUser, fmt.Errorf("failed to update user: %w", err)
	}

	// Convert models.User to type T
	updatedUser, ok := any(result).(T)
	if !ok {
		return updatedUser, fmt.Errorf("failed to convert user to type T")
	}

	return updatedUser, nil
}

// DeleteUser removes a user by their ID
func (s *userService[T]) DeleteUser(id string) error {
	ctx := context.Background()

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = s.userRepo.Delete(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// CheckUserExists verifies if a user exists by their ID
// func (s *userService[T]) CheckUserExists(id string) (bool, error) {
// 	ctx := context.Background()

// 	parsedID, err := uuid.Parse(id)
// 	if err != nil {
// 		return false, fmt.Errorf("invalid user ID: %w", err)
// 	}

// 	exists, err := s.userRepo.Exists(ctx, parsedID)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check if user exists: %w", err)
// 	}

// 	return exists, nil
// }
