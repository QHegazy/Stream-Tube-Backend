package services

import (
	"authService/internal/models"

	"github.com/google/uuid"
)

type MFAService interface {
	CreateMFA(mfa *models.MFASettings) (uuid.UUID, error)
	GetMFA(id uuid.UUID) (models.MFASettings, error)
	UpdateMFA(mfa *models.MFASettings)  error
	DeleteMFA(id string) error
	

}


