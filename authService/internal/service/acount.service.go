package services

import (
	"authService/internal/models"
	"authService/internal/repository"
	"errors"
	"log"
	"time"

	"github.com/markbates/goth"
)

type AccountService interface {
    CreateLocalAuthAccount(user *models.LocalUser, result *chan repository.ResultChan[models.User])
    CreateOAuthAccount(user goth.User, result *chan repository.ResultChan[models.User])
}

type accountService struct{}

func NewAccountService() AccountService {
    return &accountService{}
}

func (s *accountService) CreateLocalAuthAccount(user *models.LocalUser, result *chan repository.ResultChan[models.User]) {
  
}

func (s *accountService) CreateOAuthAccount(user goth.User, result *chan repository.ResultChan[models.User]) {
    userService := NewUserService[goth.User]()
    userResult := make(chan repository.ResultChan[models.User])
    go userService.CreateUser(user, &userResult) 

    select {
    case res := <-userResult:
        if res.Error != nil {
            *result <- repository.ResultChan[models.User]{Error: res.Error}
            return
        }

        oauthUser := models.OAuthUser{
            UserID:         res.Data.ID,
            Provider:       models.OAuthProvider(user.Provider),
            ProviderUserID: user.UserID,
            AccessToken:    user.AccessToken,
            RefreshToken:   user.RefreshToken,
            ExpiresAt:      user.ExpiresAt,
        }
        
        oauthResult := make(chan repository.ResultChan[models.OAuthUser])
        go NewOAuthService().CreateOAuthUser(&oauthUser, &oauthResult) 

        select {
        case oauthRes := <-oauthResult:
            if oauthRes.Error != nil {
                *result <- repository.ResultChan[models.User]{Error: oauthRes.Error}
                return
            }

            *result <- repository.ResultChan[models.User]{Data: res.Data}

        case <-time.After(10 * time.Second): 
            log.Println("Timeout creating OAuth user")
            *result <- repository.ResultChan[models.User]{Error: errors.New("timeout creating OAuth user")}
        }
    case <-time.After(10 * time.Second): 
        log.Println("Timeout creating user")
        *result <- repository.ResultChan[models.User]{Error: errors.New("timeout creating user")}
    }
}

