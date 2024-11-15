package services

import (
	"authService/internal/models"
	"authService/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type OAuthService interface {
	CreateOAuthUser(oauthUser *models.OAuthUser, result *chan repository.ResultChan[models.OAuthUser])
	GetOAuthUser(id string, result *chan repository.ResultChan[models.OAuthUser])
	UpdateOAuthUser(oauthUser *models.OAuthUser, result *chan repository.ResultChan[models.OAuthUser])
	DeleteOAuthUser(id string, result *chan repository.ResultChan[models.OAuthUser])
}

type oauthService struct {
	oauthRepo repository.OAuthRepository
}

func NewOAuthService() OAuthService {
	return &oauthService{
		oauthRepo: *repository.NewOAuthRepository(),
	}
}

func (s *oauthService) CreateOAuthUser(oauthUser *models.OAuthUser, result *chan repository.ResultChan[models.OAuthUser]) {
	ctx := context.Background()
	go func() {
		s.oauthRepo.Insert(ctx, oauthUser, result)
	}()
}

func (s *oauthService) GetOAuthUser(id string, result *chan repository.ResultChan[models.OAuthUser]) {
	ctx := context.Background()
	go func() {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			*result <- repository.ResultChan[models.OAuthUser]{Error: fmt.Errorf("invalid id: %w", err)}
			return
		}
		s.oauthRepo.Query(ctx, parsedID, result)
	}()
}

func (s *oauthService) UpdateOAuthUser(oauthUser *models.OAuthUser, result *chan repository.ResultChan[models.OAuthUser]) {
	ctx := context.Background()
	go func() {
		s.oauthRepo.Update(ctx, oauthUser, result)
	}()
}

func (s *oauthService) DeleteOAuthUser(id string, result *chan repository.ResultChan[models.OAuthUser]) {
	ctx := context.Background()
	go func() {
		parsedID, err := uuid.Parse(id)
		if err != nil {
			*result <- repository.ResultChan[models.OAuthUser]{Error: fmt.Errorf("invalid id: %w", err)}
			return
		}
		s.oauthRepo.Delete(ctx, parsedID, result)
	}()
}
