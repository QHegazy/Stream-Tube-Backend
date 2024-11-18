package services

import (
	"authService/internal/models"
	"authService/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type OAuthService interface {
	CreateOAuthUser(oauthUser *models.OAuthUser) (models.OAuthUser, error)
	GetOAuthUser(id string) (models.OAuthUser, error)
	UpdateOAuthUser(oauthUser *models.OAuthUser) (models.OAuthUser, error)
	DeleteOAuthUser(id string) error
}

type oauthService struct {
	oauthRepo repository.OAuthRepository
}

func NewOAuthService() OAuthService {
	return &oauthService{
		oauthRepo: *repository.NewOAuthRepository(),
	}
}

func (s *oauthService) CreateOAuthUser(oauthUser *models.OAuthUser) (models.OAuthUser, error) {
	ctx := context.Background()
	createdUser, err := s.oauthRepo.Insert(ctx, oauthUser)
	if err != nil {
		return models.OAuthUser{}, fmt.Errorf("failed to create OAuth user: %w", err)
	}
	return createdUser, nil
}

func (s *oauthService) GetOAuthUser(id string) (models.OAuthUser, error) {
	ctx := context.Background()

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.OAuthUser{}, fmt.Errorf("invalid user ID: %w", err)
	}

	user, err := s.oauthRepo.Query(ctx, parsedID)
	if err != nil {
		return models.OAuthUser{}, fmt.Errorf("failed to retrieve OAuth user: %w", err)
	}

	return user, nil
}

func (s *oauthService) UpdateOAuthUser(oauthUser *models.OAuthUser) (models.OAuthUser, error) {
	ctx := context.Background()

	updatedUser, err := s.oauthRepo.Update(ctx, oauthUser)
	if err != nil {
		return models.OAuthUser{}, fmt.Errorf("failed to update OAuth user: %w", err)
	}

	return updatedUser, nil
}

func (s *oauthService) DeleteOAuthUser(id string) error {
	ctx := context.Background()

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	err = s.oauthRepo.Delete(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete OAuth user: %w", err)
	}

	return nil
}
