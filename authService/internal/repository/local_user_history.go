package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PasswordHistoryRepository struct{}

func NewPasswordHistoryRepository() *PasswordHistoryRepository {
	return &PasswordHistoryRepository{}
}

func (r *PasswordHistoryRepository) Insert(ctx context.Context, passwordHistory *models.PasswordHistory) (models.PasswordHistory, error) {
	query := `
		INSERT INTO auth.password_history (user_id, password_hash, changed_at)
		VALUES ($1, $2, $3) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, passwordHistory.UserID, passwordHistory.PasswordHash, passwordHistory.CreatedAt).Scan(&passwordHistory.ID)
	if err != nil {
		return models.PasswordHistory{}, fmt.Errorf("failed to insert password history: %w", err)
	}

	return *passwordHistory, nil
}

func (r *PasswordHistoryRepository) Query(ctx context.Context, id uuid.UUID) (models.PasswordHistory, error) {
	query := `
		SELECT id, user_id, password_hash, changed_at
		FROM auth.password_history WHERE id = $1`
	pool := db.GetPools().ReadPool

	passwordHistory := models.PasswordHistory{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&passwordHistory.ID,
		&passwordHistory.UserID,
		&passwordHistory.PasswordHash,
		&passwordHistory.CreatedAt,
	)
	if err != nil {
		return models.PasswordHistory{}, fmt.Errorf("failed to query password history: %w", err)
	}

	return passwordHistory, nil
}

func (r *PasswordHistoryRepository) Update(ctx context.Context, passwordHistory *models.PasswordHistory) (models.PasswordHistory, error) {
	query := `
		UPDATE auth.password_history
		SET user_id = $1, password_hash = $2, changed_at = $3
		WHERE id = $4 RETURNING id`
	pool := db.GetPools().UpdatePool
	updatedPasswordHistory := models.PasswordHistory{}
	err := pool.QueryRow(ctx, query, passwordHistory.UserID, passwordHistory.PasswordHash, passwordHistory.CreatedAt, passwordHistory.ID).Scan(&updatedPasswordHistory.ID)
	if err != nil {
		return models.PasswordHistory{}, fmt.Errorf("failed to update password history: %w", err)
	}

	return updatedPasswordHistory, nil
}

func (r *PasswordHistoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.password_history WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete password history: %w", err)
	}

	return nil
}

func (r *PasswordHistoryRepository) QueryByUserID(ctx context.Context, UserID uuid.UUID) ([]models.PasswordHistory, error) {
	query := `SELECT id, user_id, password_hash, changed_at FROM auth.password_history WHERE user_id = $1`
	pool := db.GetPools().ReadPool

	rows, err := pool.Query(ctx, query, UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to query password history by local user ID: %w", err)
	}
	defer rows.Close()

	var passwordHistories []models.PasswordHistory
	for rows.Next() {
		var passwordHistory models.PasswordHistory
		err := rows.Scan(
			&passwordHistory.ID,
			&passwordHistory.UserID,
			&passwordHistory.PasswordHash,
			&passwordHistory.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan password history row: %w", err)
		}
		passwordHistories = append(passwordHistories, passwordHistory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over password history rows: %w", err)
	}

	return passwordHistories, nil
}

func (r *PasswordHistoryRepository) DeleteByUserID(ctx context.Context, UserID uuid.UUID) error {
	query := `DELETE FROM auth.password_history WHERE user_id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, UserID)
	if err != nil {
		return fmt.Errorf("failed to delete password history by local user ID: %w", err)
	}

	return nil
}
func (r *PasswordHistoryRepository) GeneratePasswordHistory(ctx context.Context, UserID uuid.UUID, passwordHash string) (models.PasswordHistory, error) {
	passwordHistory := models.PasswordHistory{
		UserID:  UserID,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	insertedPasswordHistory, err := r.Insert(ctx, &passwordHistory)
	if err != nil {
		return models.PasswordHistory{}, fmt.Errorf("failed to generate password history: %w", err)
	}
	return insertedPasswordHistory, nil
}
func (r *PasswordHistoryRepository) DeleteOldPasswordHistory(ctx context.Context, UserID uuid.UUID, maxHistory int) error {
	query := `DELETE FROM auth.password_history WHERE user_id = $1 AND id NOT IN (SELECT id FROM auth.password_history WHERE user_id = $1 ORDER BY changed_at DESC LIMIT $2)`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, UserID, maxHistory)
	if err != nil {
		return fmt.Errorf("failed to delete old password history: %w", err)
	}

	return nil
}
func (r *PasswordHistoryRepository) CountPasswordHistoryByUserID(ctx context.Context, UserID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM auth.password_history WHERE user_id = $1`
	pool := db.GetPools().ReadPool

	var count int
	err := pool.QueryRow(ctx, query, UserID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count password history by local user ID: %w", err)
	}

	return count, nil
}


