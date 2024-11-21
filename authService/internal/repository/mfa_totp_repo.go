package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type MFATOTPRepository struct{}

func NewMFATOTPRepository() *MFATOTPRepository {
	return &MFATOTPRepository{}
}

func (r *MFATOTPRepository) Insert(ctx context.Context, mfaTOTP *models.MFATOTP) (models.MFATOTP, error) {
	query := `
		INSERT INTO auth.mfa_totp (mfa_settings_id, secret, created_at, updated_at)
		VALUES ($1, $2, $3, $4) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, mfaTOTP.MFASettingsID, mfaTOTP.Secret, time.Now(), time.Now()).Scan(&mfaTOTP.ID)
	if err != nil {
		return models.MFATOTP{}, fmt.Errorf("failed to insert mfa totp: %w", err)
	}

	return *mfaTOTP, nil
}

func (r *MFATOTPRepository) Query(ctx context.Context, id uuid.UUID) (models.MFATOTP, error) {
	query := `
		SELECT id, mfa_settings_id, secret, created_at, updated_at
		FROM auth.mfa_totp WHERE id = $1`
	pool := db.GetPools().ReadPool

	mfaTOTP := models.MFATOTP{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&mfaTOTP.ID,
		&mfaTOTP.MFASettingsID,
		&mfaTOTP.Secret,
		&mfaTOTP.CreatedAt,
		&mfaTOTP.UpdatedAt,
	)
	if err != nil {
		return models.MFATOTP{}, fmt.Errorf("failed to query mfa totp: %w", err)
	}

	return mfaTOTP, nil
}

func (r *MFATOTPRepository) Update(ctx context.Context, mfaTOTP *models.MFATOTP) (models.MFATOTP, error) {
	query := `
		UPDATE auth.mfa_totp
		SET mfa_settings_id = $1, secret = $2, updated_at = $3
		WHERE id = $4 RETURNING id`
	pool := db.GetPools().UpdatePool

	updatedMFATOTP := models.MFATOTP{}
	err := pool.QueryRow(ctx, query, mfaTOTP.MFASettingsID, mfaTOTP.Secret, time.Now(), mfaTOTP.ID).Scan(&updatedMFATOTP.ID)
	if err != nil {
		return models.MFATOTP{}, fmt.Errorf("failed to update mfa totp: %w", err)
	}

	return updatedMFATOTP, nil
}

func (r *MFATOTPRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.mfa_totp WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete mfa totp: %w", err)
	}

	return nil
}

func (r *MFATOTPRepository) QueryByMFASettingsID(ctx context.Context, mfaSettingsID uuid.UUID) (models.MFATOTP, error) {
	query := `SELECT id, mfa_settings_id, secret, created_at, updated_at FROM auth.mfa_totp WHERE mfa_settings_id = $1`
	pool := db.GetPools().ReadPool

	mfaTOTP := models.MFATOTP{}
	err := pool.QueryRow(ctx, query, mfaSettingsID).Scan(
		&mfaTOTP.ID,
		&mfaTOTP.MFASettingsID,
		&mfaTOTP.Secret,
		&mfaTOTP.CreatedAt,
		&mfaTOTP.UpdatedAt,
	)
	if err != nil {
		return models.MFATOTP{}, fmt.Errorf("failed to query mfatotp by mfa settings ID: %w", err)
	}

	return mfaTOTP, nil
}


func (r *MFATOTPRepository) InsertBackupCode(ctx context.Context, backupCode *models.BackupCode) (models.BackupCode, error) {
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

func (r *MFATOTPRepository) QueryBackupCode(ctx context.Context, id uuid.UUID) (models.BackupCode, error) {
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

func (r *MFATOTPRepository) UpdateBackupCode(ctx context.Context, backupCode *models.BackupCode) (models.BackupCode, error) {
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

func (r *MFATOTPRepository) DeleteBackupCode(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.backup_codes WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete backup code: %w", err)
	}

	return nil
}

func (r *MFATOTPRepository) QueryBackupCodesByMFATOTP(ctx context.Context, mfaTOTPID uuid.UUID) ([]models.BackupCode, error) {
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

