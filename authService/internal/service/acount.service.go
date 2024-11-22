package services

import (
	"authService/internal/models"
	"authService/utils"
	"sync"
	"time"

	"github.com/markbates/goth"
)

// AccountService defines the interface for account-related operations
type AccountService[T any] interface {
	CreateLocalAuthAccount(user *models.LocalUser,s *sync.WaitGroup) (T, error)
	CreateOAuthAccount(user goth.User,s *sync.WaitGroup) (T, error)
}

type accountService[T any] struct {
}

func NewAccountService[T any]() AccountService[T] {
	return &accountService[T]{}
}

func (s *accountService[T]) CreateLocalAuthAccount(user *models.LocalUser,ss *sync.WaitGroup) (T, error) {
	defer ss.Done()
	var result T
	return result, nil
}

func (s *accountService[T]) CreateOAuthAccount(user goth.User,ss *sync.WaitGroup) (T, error) {
	defer ss.Done()
	var result T
    now := time.Now().UTC()
	// Create a new user via the userService
	userService := NewUserService[T]()
	username := user.Name + "_" + utils.GenerateRandomString(4)
	newUser := models.User{
		Username:      username,
		Email:         user.Email,
		AuthMethod:    models.AuthMethodOAuth,
		Status:        "active",
		EmailVerified: now,
		LastActiveAt:  now,
	}
	userResult, err := userService.CreateUser(&newUser)
	if err != nil {
		return result, err
	}

	// Create OAuth user data
	oauthUser := models.OAuthUser{
		UserID:         userResult,
		Provider:       models.OAuthProvider(user.Provider),
		ProviderUserID: user.UserID,
		AccessToken:    user.AccessToken,
		RefreshToken:   user.RefreshToken,
		ExpiresAt:      user.ExpiresAt,
	}

	// Create OAuth user in the database
	oauthService := NewOAuthService()
	_, err = oauthService.CreateOAuthUser(&oauthUser)
	if err != nil {
		return result, err
	}

	// Convert userResult to type T if needed
	// This depends on your specific implementation
	
	return result, nil
}