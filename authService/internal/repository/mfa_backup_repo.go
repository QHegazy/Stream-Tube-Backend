package repository

import (
	"authService/db"
	"authService/internal/models"
	"authService/utils"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MFABackupCodesRepository struct{}

func NewMFABackupCodesRepository() *MFABackupCodesRepository {
	return &MFABackupCodesRepository{}
}

func (r *MFABackupCodesRepository) GenerateCodes(ctx context.Context, mfaTotpID uuid.UUID, count int) ([]models.BackupCode, error) {
	codes := make([]models.BackupCode, count)
	for i := 0; i < count; i++ {
		code := models.BackupCode{
			MFATOTPID:  mfaTotpID,
			Code:       utils.GenerateRandomString(6),
			UsedCount:  0,
			CreatedAt: time.Now(),
		}
		insertedCode, err := r.Insert(ctx, &code)
		if err != nil {
			return nil, fmt.Errorf("failed to insert backup code: %w", err)
		}
		codes[i] = insertedCode
	}
	return codes, nil
}

func (r *MFABackupCodesRepository) Insert(ctx context.Context, backupCode *models.BackupCode) (models.BackupCode, error) {
	query := `
		INSERT INTO auth.backup_codes (mfa_totpid, code, used_count, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, backupCode.MFATOTPID, backupCode.Code, backupCode.UsedCount, time.Now()).Scan(&backupCode.ID)
	if err != nil {
		return models.BackupCode{}, fmt.Errorf("failed to insert backup code: %w", err)
	}

	return *backupCode, nil
}

func (r *MFABackupCodesRepository) Query(ctx context.Context, id uuid.UUID) (models.BackupCode, error) {
	query := `
		SELECT id, mfa_totpid, code, used_count, created_at
		FROM auth.backup_codes WHERE id = $1`
	pool := db.GetPools().ReadPool

	backupCode := models.BackupCode{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&backupCode.ID,
		&backupCode.MFATOTPID,
		&backupCode.Code,
		&backupCode.UsedCount,
		&backupCode.CreatedAt,
	)
	if err != nil {
		return models.BackupCode{}, fmt.Errorf("failed to query backup code: %w", err)
	}

	return backupCode, nil
}

func (r *MFABackupCodesRepository) Update(ctx context.Context, backupCode *models.BackupCode) (models.BackupCode, error) {
	query := `
		UPDATE auth.backup_codes
		SET mfa_totpid = $1, code = $2, used_count = $3, created_at = $4
		WHERE id = $5 RETURNING id`
	pool := db.GetPools().UpdatePool

	updatedBackupCode := models.BackupCode{}
	err := pool.QueryRow(ctx, query, backupCode.MFATOTPID, backupCode.Code, backupCode.UsedCount, backupCode.CreatedAt, backupCode.ID).Scan(&updatedBackupCode.ID)
	if err != nil {
		return models.BackupCode{}, fmt.Errorf("failed to update backup code: %w", err)
	}

	return updatedBackupCode, nil
}

func (r *MFABackupCodesRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.backup_codes WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete backup code: %w", err)
	}

	return nil
}

func (r *MFABackupCodesRepository) QueryByMFATOTP(ctx context.Context, mfaTOTPID uuid.UUID) ([]models.BackupCode, error) {
	query := `SELECT id, mfa_totpid, code, used_count, created_at FROM auth.backup_codes WHERE mfa_totpid = $1`
	pool := db.GetPools().ReadPool

	rows, err := pool.Query(ctx, query, mfaTOTPID)
	if err != nil {
		return nil, fmt.Errorf("failed to query backup codes by mfa totp ID: %w", err)
	}
	defer rows.Close()

	var backupCodes []models.BackupCode
	for rows.Next() {
		var backupCode models.BackupCode
		err := rows.Scan(
			&backupCode.ID,
			&backupCode.MFATOTPID,
			&backupCode.Code,
			&backupCode.UsedCount,
			&backupCode.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan backup code row: %w", err)
		}
		backupCodes = append(backupCodes, backupCode)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over backup code rows: %w", err)
	}

	return backupCodes, nil
}

func (r *MFABackupCodesRepository) UseBackupCode(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE auth.backup_codes SET used_count = used_count + 1 WHERE id = $1`
	pool := db.GetPools().UpdatePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to use backup code: %w", err)
	}

	return nil
}

func (r *MFABackupCodesRepository) DeleteByMFATOTP(ctx context.Context, mfaTotpID uuid.UUID) error {
	query := `DELETE FROM auth.backup_codes WHERE mfa_totpid = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, mfaTotpID)
	if err != nil {
		return fmt.Errorf("failed to delete backup codes by mfa totp ID: %w", err)
	}
	return nil
}
