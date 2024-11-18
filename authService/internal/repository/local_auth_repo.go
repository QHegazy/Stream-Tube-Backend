package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type LocalAuthRepository struct{}

func NewLocalAuthRepository() *LocalAuthRepository {
	return &LocalAuthRepository{}
}

func (r *LocalAuthRepository) Insert(ctx context.Context, localUser models.LocalUser) (models.LocalUser, error) {
	query := `
		INSERT INTO auth.local_users 
		(user_id, password, last_password_change, password_history, force_password_change, password_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query,
		localUser.UserID, localUser.Password, localUser.LastPasswordChange,
		localUser.PasswordHistory, localUser.ForcePasswordChange, localUser.PasswordExpiresAt,
	).Scan(&localUser.ID)
	if err != nil {
		return models.LocalUser{}, fmt.Errorf("failed to insert local user: %w", err)
	}

	return localUser, nil
}

func (r *LocalAuthRepository) Update(ctx context.Context, localUser models.LocalUser) (models.LocalUser, error) {
	query := `
		UPDATE auth.local_users 
		SET user_id = $1, password = $2, last_password_change = $3,
		password_history = $4, force_password_change = $5, password_expires_at = $6
		WHERE id = $7 RETURNING id`
	pool := db.GetPools().UpdatePool

	err := pool.QueryRow(ctx, query,
		localUser.UserID, localUser.Password, localUser.LastPasswordChange,
		localUser.PasswordHistory, localUser.ForcePasswordChange, localUser.PasswordExpiresAt, localUser.ID,
	).Scan(&localUser.ID)
	if err != nil {
		return models.LocalUser{}, fmt.Errorf("failed to update local user: %w", err)
	}

	return localUser, nil
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

func (r *LocalAuthRepository) Query(ctx context.Context, id uuid.UUID) (models.LocalUser, error) {
	query := `
		SELECT id, user_id, password, last_password_change, password_history, force_password_change, password_expires_at
		FROM auth.local_users WHERE id = $1`
	pool := db.GetPools().ReadPool

	localUser := models.LocalUser{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&localUser.ID, &localUser.UserID, &localUser.Password, &localUser.LastPasswordChange,
		&localUser.PasswordHistory, &localUser.ForcePasswordChange, &localUser.PasswordExpiresAt,
	)
	if err != nil {
		return models.LocalUser{}, fmt.Errorf("failed to query local user: %w", err)
	}

	return localUser, nil
}
