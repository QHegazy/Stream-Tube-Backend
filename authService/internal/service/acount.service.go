package services

import (
	dto "authService/Dto"
	"authService/config"
	"authService/internal/models"
	"authService/utils"
	"fmt"
	"sync"
	"time"

	"github.com/markbates/goth"
)

// AccountService defines the interface for account-related operations
type AccountService[T any] interface {
	CreateLocalAuthAccount(user *dto.RegisterLocalUser,s *sync.WaitGroup) (T, error)
	CreateOAuthAccount(user goth.User,s *sync.WaitGroup) (T, error)
}

type accountService[T any] struct {
}

func NewAccountService[T any]() AccountService[T] {
	return &accountService[T]{}
}

func (s *accountService[T]) CreateLocalAuthAccount(user *dto.RegisterLocalUser, wg *sync.WaitGroup) (T, error) {
	defer wg.Done()

	var result T
	userService := NewUserService[T]()
	newUser := models.User{
		Username:      user.Username,
		Email:         user.Email,
		AuthMethod:    models.AuthMethodLocal,
		Status:        "active",
		EmailVerified: time.Time{},
		LastActiveAt:  time.Now().UTC(),
	}

	userResult, err := userService.CreateUser(&newUser)
	if err != nil {
		return result, err
	}
	PasswordHash,err := utils.HashPasswordWithSalt(user.Password,config.GetConfig().Salt)
	if  err !=nil {
		return result ,err
	}
	
	localUserService := NewLocalUserService()
	newLocalUser := models.LocalUser{
		UserID:       userResult,
		PasswordHash: PasswordHash,
	}
	err = localUserService.CreateLocalUser(&newLocalUser)
	if err != nil {
		return result, err
	}

	// Convert userResult to T
	if convertedResult, ok := any(userResult).(T); ok {
		return convertedResult, nil
	}

	return result, fmt.Errorf("type mismatch: cannot convert uuid.UUID to %T", result)
}


func (s *accountService[T]) CreateOAuthAccount(user goth.User,wg *sync.WaitGroup) (T, error) {
	defer wg.Done()
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


	return result, nil
}

