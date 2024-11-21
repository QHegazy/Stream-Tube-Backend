package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MFASettingsRepository struct{}

func NewMFASettingsRepository() *MFASettingsRepository {
	return &MFASettingsRepository{}
}

func (r *MFASettingsRepository) Insert(ctx context.Context, mfaSettings *models.MFASettings) (models.MFASettings, error) {
	query := `
		INSERT INTO auth.mfa_settings (user_id, mfa_type, last_mfa_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, mfaSettings.UserID, mfaSettings.MFAType, mfaSettings.LastMFAAt, time.Now(), time.Now()).Scan(&mfaSettings.ID)
	if err != nil {
		return models.MFASettings{}, fmt.Errorf("failed to insert mfa settings: %w", err)
	}

	return *mfaSettings, nil
}

func (r *MFASettingsRepository) Query(ctx context.Context, id uuid.UUID) (models.MFASettings, error) {
	query := `
		SELECT id, user_id, mfa_type, last_mfa_at, created_at, updated_at
		FROM auth.mfa_settings WHERE id = $1`
	pool := db.GetPools().ReadPool

	mfaSettings := models.MFASettings{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&mfaSettings.ID,
		&mfaSettings.UserID,
		&mfaSettings.MFAType,
		&mfaSettings.LastMFAAt,
		&mfaSettings.CreatedAt,
		&mfaSettings.UpdatedAt,
	)
	if err != nil {
		return models.MFASettings{}, fmt.Errorf("failed to query mfa settings: %w", err)
	}

	return mfaSettings, nil
}

func (r *MFASettingsRepository) Update(ctx context.Context, mfaSettings *models.MFASettings) (models.MFASettings, error) {
	query := `
		UPDATE auth.mfa_settings
		SET user_id = $1, mfa_type = $2, last_mfa_at = $3, updated_at = $4
		WHERE id = $5 RETURNING id`
	pool := db.GetPools().UpdatePool

	updatedMFASettings := models.MFASettings{}
	err := pool.QueryRow(ctx, query, mfaSettings.UserID, mfaSettings.MFAType, mfaSettings.LastMFAAt, time.Now(), mfaSettings.ID).Scan(&updatedMFASettings.ID)
	if err != nil {
		return models.MFASettings{}, fmt.Errorf("failed to update mfa settings: %w", err)
	}

	return updatedMFASettings, nil
}

func (r *MFASettingsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.mfa_settings WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete mfa settings: %w", err)
	}

	return nil
}

func (r *MFASettingsRepository) QueryByUser(ctx context.Context, userID uuid.UUID) (models.MFASettings, error) {
	query := `SELECT id, user_id, mfa_type, last_mfa_at, created_at, updated_at FROM auth.mfa_settings WHERE user_id = $1`
	pool := db.GetPools().ReadPool

	mfaSettings := models.MFASettings{}
	err := pool.QueryRow(ctx, query, userID).Scan(
		&mfaSettings.ID,
		&mfaSettings.UserID,
		&mfaSettings.MFAType,
		&mfaSettings.LastMFAAt,
		&mfaSettings.CreatedAt,
		&mfaSettings.UpdatedAt,
	)
	if err != nil {
		return models.MFASettings{}, fmt.Errorf("failed to query mfa settings by user ID: %w", err)
	}

	return mfaSettings, nil
}

