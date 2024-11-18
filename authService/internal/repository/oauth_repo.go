package repository

import (
	"authService/config"
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type OAuthRepository struct{}

func NewOAuthRepository() *OAuthRepository {
	return &OAuthRepository{}
}

func (r *OAuthRepository) Insert(ctx context.Context, oauthUser *models.OAuthUser) (models.OAuthUser, error) {
	db.InitDB(config.GetDBConfig())
	query := `
		INSERT INTO auth.oauth_users (user_id, provider, provider_user_id, access_token, refresh_token, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, oauthUser.UserID, oauthUser.Provider, oauthUser.ProviderUserID, oauthUser.AccessToken, oauthUser.RefreshToken, oauthUser.ExpiresAt).Scan(&oauthUser.ID)
	if err != nil {
		return models.OAuthUser{}, fmt.Errorf("failed to insert oauth user: %w", err)
	}

	return *oauthUser, nil
}

func (r *OAuthRepository) Query(ctx context.Context, id uuid.UUID) (models.OAuthUser, error) {
	db.InitDB(config.GetDBConfig())
	query := `
		SELECT id, user_id, provider, provider_user_id, access_token, refresh_token, expires_at
		FROM auth.oauth_users WHERE id = $1`
	pool := db.GetPools().ReadPool

	oauthUser := models.OAuthUser{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&oauthUser.ID,
		&oauthUser.UserID,
		&oauthUser.Provider,
		&oauthUser.ProviderUserID,
		&oauthUser.AccessToken,
		&oauthUser.RefreshToken,
		&oauthUser.ExpiresAt,
	)
	if err != nil {
		return models.OAuthUser{}, fmt.Errorf("failed to query oauth user: %w", err)
	}

	return oauthUser, nil
}

func (r *OAuthRepository) Update(ctx context.Context, oauthUser *models.OAuthUser) (models.OAuthUser, error) {
	db.InitDB(config.GetDBConfig())
	query := `
		UPDATE auth.oauth_users
		SET user_id = $1, provider = $2, provider_user_id = $3, access_token = $4, refresh_token = $5, expires_at = $6
		WHERE id = $7 RETURNING id`
	pool := db.GetPools().UpdatePool

	updatedOAuthUser := models.OAuthUser{}
	err := pool.QueryRow(ctx, query, oauthUser.UserID, oauthUser.Provider, oauthUser.ProviderUserID, oauthUser.AccessToken, oauthUser.RefreshToken, oauthUser.ExpiresAt, oauthUser.ID).Scan(&updatedOAuthUser.ID)
	if err != nil {
		return models.OAuthUser{}, fmt.Errorf("failed to update oauth user: %w", err)
	}

	return updatedOAuthUser, nil
}

func (r *OAuthRepository) Delete(ctx context.Context, id uuid.UUID)  error {
	db.InitDB(config.GetDBConfig())
	query := `DELETE FROM auth.oauth_users WHERE id = $1 RETURNING id`
	pool := db.GetPools().DeletePool

	deletedOAuthUser := models.OAuthUser{}
	err := pool.QueryRow(ctx, query, id).Scan(&deletedOAuthUser.ID)
	if err != nil {
		return fmt.Errorf("failed to delete oauth user: %w", err)
	}

	return nil
}
