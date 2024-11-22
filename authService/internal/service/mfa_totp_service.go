package services

import (
	"authService/internal/models"
	"authService/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type MFATOTPService interface {
	CreateMFATOTP(mfaTOTP *models.MFATOTP) (uuid.UUID, error)
	GetMFATOTP(id string) (models.MFATOTP, error)
	UpdateMFATOTP(mfaTOTP *models.MFATOTP) (models.MFATOTP, error)
	DeleteMFATOTP(id string) error
	GetMFATOTPByMFASettingsID(mfaSettingsID string) (models.MFATOTP, error)
	GenerateBackupCodes(mfaTotpID uuid.UUID, count int) ([]models.BackupCode, error)
	UseBackupCode(id string) error
	DeleteBackupCodesByMFATOTP(mfaTotpID string) error
}

type mfaTOTPService struct {
	mfaTOTPRepo      repository.MFATOTPRepository
	backupCodeRepo   repository.MFABackupCodesRepository
}

func NewMFATOTPService() MFATOTPService {
	return &mfaTOTPService{
		mfaTOTPRepo: *repository.NewMFATOTPRepository(),
		backupCodeRepo: *repository.NewMFABackupCodesRepository(),
	}
}

func (s *mfaTOTPService) CreateMFATOTP(mfaTOTP *models.MFATOTP) (uuid.UUID, error) {
	ctx := context.Background()
	createdMFATOTP, err := s.mfaTOTPRepo.Insert(ctx, mfaTOTP)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create mfa totp: %w", err)
	}
	return createdMFATOTP.ID, nil
}

func (s *mfaTOTPService) GetMFATOTP(id string) (models.MFATOTP, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.MFATOTP{}, fmt.Errorf("invalid mfa totp ID: %w", err)
	}
	mfaTOTP, err := s.mfaTOTPRepo.Query(ctx, parsedID)
	if err != nil {
		return models.MFATOTP{}, fmt.Errorf("failed to get mfa totp: %w", err)
	}
	return mfaTOTP, nil
}

func (s *mfaTOTPService) UpdateMFATOTP(mfaTOTP *models.MFATOTP) (models.MFATOTP, error) {
	ctx := context.Background()
	updatedMFATOTP, err := s.mfaTOTPRepo.Update(ctx, mfaTOTP)
	if err != nil {
		return models.MFATOTP{}, fmt.Errorf("failed to update mfa totp: %w", err)
	}
	return updatedMFATOTP, nil
}

func (s *mfaTOTPService) DeleteMFATOTP(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid mfa totp ID: %w", err)
	}
	err = s.mfaTOTPRepo.Delete(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete mfa totp: %w", err)
	}
	return nil
}

func (s *mfaTOTPService) GetMFATOTPByMFASettingsID(mfaSettingsID string) (models.MFATOTP, error) {
	ctx := context.Background()
	parsedMFASettingsID, err := uuid.Parse(mfaSettingsID)
	if err != nil {
		return models.MFATOTP{}, fmt.Errorf("invalid mfa settings ID: %w", err)
	}
	mfaTOTP, err := s.mfaTOTPRepo.QueryByMFASettingsID(ctx, parsedMFASettingsID)
	if err != nil {
		return models.MFATOTP{}, fmt.Errorf("failed to get mfa totp by mfa settings ID: %w", err)
	}
	return mfaTOTP, nil
}

func (s *mfaTOTPService) GenerateBackupCodes(mfaTotpID uuid.UUID, count int) ([]models.BackupCode, error) {
	ctx := context.Background()
	backupCodes, err := s.backupCodeRepo.GenerateCodes(ctx, mfaTotpID, count)
	if err != nil {
		return nil, fmt.Errorf("failed to generate backup codes: %w", err)
	}
	return backupCodes, nil
}

func (s *mfaTOTPService) UseBackupCode(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid backup code ID: %w", err)
	}
	err = s.backupCodeRepo.UseBackupCode(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to use backup code: %w", err)
	}
	return nil
}

func (s *mfaTOTPService) DeleteBackupCodesByMFATOTP(mfaTotpID string) error {
	ctx := context.Background()
	parsedMFATOTP, err := uuid.Parse(mfaTotpID)
	if err != nil {
		return fmt.Errorf("invalid mfa totp ID: %w", err)
	}
	err = s.backupCodeRepo.DeleteByMFATOTP(ctx, parsedMFATOTP)
	if err != nil {
		return fmt.Errorf("failed to delete backup codes by mfa totp ID: %w", err)
	}
	return nil
}

// func Totp(username string){
// 	secret, err := totp.Generate(totp.GenerateOpts{
// 		Issuer:      "Stream Tube",
// 		AccountName: username,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	qr, err := qrcode.Encode(secret.URL(), qrcode.Medium, 256)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(secret.URL())
// 	fmt.Println(string(qr))

// }

// func TotpVerfaiy(code string,secret string){
// 	valid := totp.Validate(code,secret )
// 	fmt.Println(valid)

// }