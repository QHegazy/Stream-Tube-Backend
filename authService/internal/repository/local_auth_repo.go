package repository

import (
	dto "authService/Dto"
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type LocalAuthRepository struct{}

func NewLocalAuthRepository() *LocalAuthRepository {
	return &LocalAuthRepository{}
}

func (r *LocalAuthRepository) Insert(ctx context.Context, localUser *models.LocalUser) (models.LocalUser, error) {
	query := `
		INSERT INTO auth.local_users (user_id, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, localUser.UserID, localUser.PasswordHash, time.Now(), time.Now()).Scan(&localUser.ID)
	if err != nil {
		return models.LocalUser{}, fmt.Errorf("failed to insert local user: %w", err)
	}

	return *localUser, nil
}

func (r *LocalAuthRepository) Query(ctx context.Context, localUser models.LocalUser) (models.LocalUser, error) {
	query := `
		SELECT id, user_id, password_hash, created_at, updated_at
		FROM auth.local_users WHERE password_hash = $1 AND user_id = $2`
	pool := db.GetPools().ReadPool
	err := pool.QueryRow(ctx, query,&localUser.PasswordHash,&localUser.UserID).Scan(
		&localUser.ID,
		&localUser.UserID,
		&localUser.PasswordHash,
		&localUser.CreatedAt,
		&localUser.UpdatedAt,
	)
	if err != nil {
		return models.LocalUser{}, fmt.Errorf("failed to query local user: %w", err)
	}

	return localUser, nil
}

func (r *LocalAuthRepository) Update(ctx context.Context, localUser *models.LocalUser) (models.LocalUser, error) {
	query := `
		UPDATE auth.local_users
		SET user_id = $1, password_hash = $2, updated_at = $3
		WHERE id = $4 RETURNING id`
	pool := db.GetPools().UpdatePool

	updatedLocalUser := models.LocalUser{}
	err := pool.QueryRow(ctx, query, localUser.UserID, localUser.PasswordHash, time.Now(), localUser.ID).Scan(&updatedLocalUser.ID)
	if err != nil {
		return models.LocalUser{}, fmt.Errorf("failed to update local user: %w", err)
	}

	return updatedLocalUser, nil
}

func (r *LocalAuthRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.local_users WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete local user: %w", err)
	}

	return nil
}

func (r *LocalAuthRepository) QueryByUser(ctx context.Context, userID uuid.UUID) (models.LocalUser, error) {
	query := `SELECT id, user_id, password_hash, created_at, updated_at FROM auth.local_users WHERE user_id = $1`
	pool := db.GetPools().ReadPool

	localUser := models.LocalUser{}
	err := pool.QueryRow(ctx, query, userID).Scan(
		&localUser.ID,
		&localUser.UserID,
		&localUser.PasswordHash,
		&localUser.CreatedAt,
		&localUser.UpdatedAt,
	)
	if err != nil {
		return models.LocalUser{}, fmt.Errorf("failed to query local user by user ID: %w", err)
	}

	return localUser, nil
}

func (r *LocalAuthRepository) UpdateLocalUserPassword(ctx context.Context, oldPasswordHash dto.LocalUserDto, newPasswordHash string) error {
	query := `UPDATE auth.local_users SET password_hash = $1, updated_at = $2 WHERE user_id = $3 AND password_hash = $4`
	pool := db.GetPools().UpdatePool

	_, err := pool.Exec(ctx, query, newPasswordHash, time.Now(), oldPasswordHash.LocalUserID, oldPasswordHash.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to update local user password: %w", err)
	}

	return nil
}

func (r *LocalAuthRepository) CheckLocalUserExists(ctx context.Context, password string) (bool, error) {
	query := `SELECT COUNT(*) FROM auth.local_users WHERE password_hash = $1`
	pool := db.GetPools().ReadPool

	var count int
	err := pool.QueryRow(ctx, query, password).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check local user exists: %w", err)
	}

	return count > 0, nil
}
