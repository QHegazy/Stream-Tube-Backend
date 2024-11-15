package services

import (
	"authService/internal/models"
	"authService/internal/repository"
	"authService/utils"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/markbates/goth"
)

// UserService interface defines methods for user management
type UserService[T any] interface {
	CreateUser(user T, result *chan repository.ResultChan[models.User])
	GetUser(id string, result *chan repository.ResultChan[models.User])
	UpdateUser(user *models.User, result *chan repository.ResultChan[models.User])
	DeleteUser(id string, result *chan repository.ResultChan[models.User])
	CheckUserExists(id string, result chan repository.ResultChan[bool])
}

// userService struct that uses a generic type T
type userService[T any] struct {
	userRepo repository.UserRepository
}

// NewUserService initializes a new UserService instance with a generic type T
func NewUserService[T any]() UserService[T] {
	return &userService[T]{
		userRepo: *repository.NewUserRepository(),
	}
}

// CreateUser inserts a new user based on Goth user information
func (s *userService[T]) CreateUser(user T, result *chan repository.ResultChan[models.User]) {
	now := time.Now().UTC()
	gothUser, ok := any(user).(goth.User)
	fmt.Println("user:",user)
	if !ok {
		*result <- repository.ResultChan[models.User]{Error: fmt.Errorf("invalid user type")}
		return
	}
	fmt.Println("gothUser:",gothUser)
	
	username := gothUser.Name + "_" + utils.GenerateRandomString(4)
	newUser := models.User{
		Username:      username,
		Email:         gothUser.Email,
		AuthMethod:    models.AuthMethodOAuth,
		Status:        "active",
		LastActiveAt:  &now,
		EmailVerified: &now,
	}
	ctx := context.Background()
	go func() {
		s.userRepo.Insert(ctx, &newUser, result)
		
	}()
}

// GetUser retrieves a user by their ID
func (s *userService[T]) GetUser(id string, result *chan repository.ResultChan[models.User]) {
	ctx := context.Background()
	go func() {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			*result <- repository.ResultChan[models.User]{Error: err}
			return
		}
		s.userRepo.Query(ctx, parsedID, result)
	}()
}

// UpdateUser updates the user information
func (s *userService[T]) UpdateUser(user *models.User, result *chan repository.ResultChan[models.User]) {
	ctx := context.Background()
	go func() {
		s.userRepo.Update(ctx, user, result)
	}()
}

// DeleteUser removes a user by their ID
func (s *userService[T]) DeleteUser(id string, result *chan repository.ResultChan[models.User]) {
	ctx := context.Background()
	go func() {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			*result <- repository.ResultChan[models.User]{Error: err}
			return
		}
		s.userRepo.Delete(ctx, parsedID, result)
	}()
}

// CheckUserExists verifies if a user exists by their ID
func (s *userService[T]) CheckUserExists(id string, result chan repository.ResultChan[bool]) {
	ctx := context.Background()
	go func() {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			result <- repository.ResultChan[bool]{Data: false, Error: err}
			return
		}
		exists := s.userRepo.Exists(ctx, parsedID)
		result <- repository.ResultChan[bool]{Data: exists, Error: nil}
	}()
}
