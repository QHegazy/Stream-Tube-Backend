package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type LocalUserFlagRepository struct{}

func NewLocalUserFlagRepository() *LocalUserFlagRepository {
	return &LocalUserFlagRepository{}
}

func (r *LocalUserFlagRepository) Insert(ctx context.Context, localUserFlag *models.LocalUserFlag) (uuid.UUID, error) {
	query := `
		INSERT INTO auth.local_user_flags (local_user_id, force_password_change, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, localUserFlag.LocalUserID, localUserFlag.ForcePasswordChange, time.Now(), time.Now()).Scan(&localUserFlag.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert local user flag: %w", err)
	}

	return localUserFlag.ID, nil
}

func (r *LocalUserFlagRepository) Query(ctx context.Context, id uuid.UUID) (models.LocalUserFlag, error) {
	query := `
		SELECT id, local_user_id, force_password_change, created_at, updated_at
		FROM auth.local_user_flags WHERE id = $1`
	pool := db.GetPools().ReadPool

	localUserFlag := models.LocalUserFlag{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&localUserFlag.ID,
		&localUserFlag.LocalUserID,
		&localUserFlag.ForcePasswordChange,
		&localUserFlag.CreatedAt,
		&localUserFlag.UpdatedAt,
	)
	if err != nil {
		return models.LocalUserFlag{}, fmt.Errorf("failed to query local user flag: %w", err)
	}

	return localUserFlag, nil
}

func (r *LocalUserFlagRepository) Update(ctx context.Context, localUserFlag *models.LocalUserFlag) (models.LocalUserFlag, error) {
	query := `
		UPDATE auth.local_user_flags
		SET local_user_id = $1, force_password_change = $2, updated_at = $3
		WHERE id = $4 RETURNING id`
	pool := db.GetPools().UpdatePool

	updatedLocalUserFlag := models.LocalUserFlag{}
	err := pool.QueryRow(ctx, query, localUserFlag.LocalUserID, localUserFlag.ForcePasswordChange, time.Now(), localUserFlag.ID).Scan(&updatedLocalUserFlag.ID)
	if err != nil {
		return models.LocalUserFlag{}, fmt.Errorf("failed to update local user flag: %w", err)
	}

	return updatedLocalUserFlag, nil
}

func (r *LocalUserFlagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.local_user_flags WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete local user flag: %w", err)
	}

	return nil
}

func (r *LocalUserFlagRepository) QueryByUser(ctx context.Context, userID uuid.UUID) (models.LocalUserFlag, error) {
	query := `SELECT id, local_user_id, force_password_change, created_at, updated_at FROM auth.local_user_flags WHERE local_user_id = $1`
	pool := db.GetPools().ReadPool

	localUserFlag := models.LocalUserFlag{}
	err := pool.QueryRow(ctx, query, userID).Scan(
		&localUserFlag.ID,
		&localUserFlag.LocalUserID,
		&localUserFlag.ForcePasswordChange,
		&localUserFlag.CreatedAt,
		&localUserFlag.UpdatedAt,
	)
	if err != nil {
		return models.LocalUserFlag{}, fmt.Errorf("failed to query local user flag by user ID: %w", err)
	}

	return localUserFlag, nil
}
func (r *LocalUserFlagRepository) GenerateLocalUserFlags(ctx context.Context, localUserID uuid.UUID, count int, forcePasswordChange bool) ([]models.LocalUserFlag, error) {
	flags := make([]models.LocalUserFlag, count)
	for i := 0; i < count; i++ {
		flag := models.LocalUserFlag{
			LocalUserID:        localUserID,
			ForcePasswordChange: forcePasswordChange,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		_, err := r.Insert(ctx, &flag)
		if err != nil {
			return nil, fmt.Errorf("failed to insert local user flag: %w", err)
		}
		flags[i] = flag
	}
	return flags, nil
}
func (r *LocalUserFlagRepository) DeleteByLocalUserID(ctx context.Context, localUserID uuid.UUID) error {
	query := `DELETE FROM auth.local_user_flags WHERE local_user_id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, localUserID)
	if err != nil {
		return fmt.Errorf("failed to delete local user flags by local user ID: %w", err)
	}

	return nil
}

