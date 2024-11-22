package services

import (
	dto "authService/Dto"
	"authService/internal/models"
	"authService/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type LocalUserService interface {
	CreateLocalUser(localUser *models.LocalUser)  error
	GetLocalUser(localUser *models.LocalUser)  error
	UpdateLocalUser(localUser *models.LocalUser) error
	DeleteLocalUser(localUser *models.LocalUser) error
	GetLocalUserByUserID(userID uuid.UUID) (models.LocalUser, error)
	CheckLocalUserExists(password string) (bool, error)
	UpdateLocalUserPassword(oldPasswordHash dto.LocalUserDto, newPasswordHash string) error
	GetLocalUserPasswordHistory(id uuid.UUID) ([]models.PasswordHistory, error)
	AddLocalUserPasswordHistory(id uuid.UUID, passwordHash string) error
	DeleteLocalUserPasswordHistory(id uuid.UUID) error
	CreateLocalUserFlag(flag *models.LocalUserFlag) (uuid.UUID, error)
	GetLocalUserFlag(id string) (models.LocalUserFlag, error)
	UpdateLocalUserFlag(flag *models.LocalUserFlag) (models.LocalUserFlag, error)
	DeleteLocalUserFlag(id string) error
	GetLocalUserFlagByUserID(userID string) (models.LocalUserFlag, error)
	CreateUserSecurityQuestion(question *models.UserSecurityQuestion)  error
	GetUserSecurityQuestion(id string) (models.UserSecurityQuestion, error)
	UpdateUserSecurityQuestion(question *models.UserSecurityQuestion) (models.UserSecurityQuestion, error)
	DeleteUserSecurityQuestion(id string) error
	GetUserSecurityQuestionByUser(userID string) ([]models.UserSecurityQuestion, error)
	DeleteUserSecurityQuestionsByUser(userID string) error
}
type localUserService struct {
	localUserRepo          repository.LocalAuthRepository
	passwordHistoryRepo    repository.PasswordHistoryRepository
	userSecurityQuestionRepo repository.UserSecurityQuestionRepository
	localUserFlagRepo      repository.LocalUserFlagRepository
}

func NewLocalUserService() LocalUserService {
	return &localUserService{
		localUserRepo:          *repository.NewLocalAuthRepository(),
		passwordHistoryRepo:    *repository.NewPasswordHistoryRepository(),
		userSecurityQuestionRepo: *repository.NewUserSecurityQuestionRepository(),
		localUserFlagRepo:      *repository.NewLocalUserFlagRepository(),
	}
}

func (s *localUserService) CreateLocalUser(localUser *models.LocalUser) error {
	ctx := context.Background()
	_, err := s.localUserRepo.Insert(ctx, localUser)
	if err != nil {
		return fmt.Errorf("failed to create local user: %w", err)
	}
	return nil
}

func (s *localUserService) GetLocalUser(localUser *models.LocalUser) error {
	ctx := context.Background()
	_, err := s.localUserRepo.Query(ctx, *localUser)
	if err != nil {
		return fmt.Errorf("failed to get local user: %w", err)
	}
	return nil
}

func (s *localUserService) UpdateLocalUser(localUser *models.LocalUser) error {
	ctx := context.Background()
	_, err := s.localUserRepo.Update(ctx, localUser)
	if err != nil {
		return fmt.Errorf("failed to update local user: %w", err)
	}
	return nil
}

func (s *localUserService) DeleteLocalUser(localUser *models.LocalUser) error {
	ctx := context.Background()
	err := s.localUserRepo.Delete(ctx, localUser.ID)
	if err != nil {
		return fmt.Errorf("failed to delete local user: %w", err)
	}
	return nil
}

func (s *localUserService) GetLocalUserByUserID(userID uuid.UUID) (models.LocalUser, error) {
	ctx := context.Background()
	localUser, err := s.localUserRepo.QueryByUser(ctx, userID)
	if err != nil {
		return models.LocalUser{}, fmt.Errorf("failed to get local user by user ID: %w", err)
	}
	return localUser, nil
}
func (s *localUserService) CheckLocalUserExists(password string) (bool, error) {
	ctx := context.Background()
	exists, err := s.localUserRepo.CheckLocalUserExists(ctx, password)
	if err != nil {
		return false, fmt.Errorf("failed to check local user exists: %w", err)
	}
	return exists, nil
}

func (s *localUserService) UpdateLocalUserPassword(oldPasswordHash dto.LocalUserDto, newPasswordHash string) error {
	ctx := context.Background()
	err := s.localUserRepo.UpdateLocalUserPassword(ctx, oldPasswordHash, newPasswordHash)
	if err != nil {
		return fmt.Errorf("failed to update local user password: %w", err)
	}
	return nil
}

func (s *localUserService) GetLocalUserPasswordHistory(id uuid.UUID) ([]models.PasswordHistory, error) {
	ctx := context.Background()
	passwordHistories, err := s.passwordHistoryRepo.QueryByUserID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get local user password history: %w", err)
	}
	return passwordHistories, nil
}

func (s *localUserService) AddLocalUserPasswordHistory(id uuid.UUID, passwordHash string) error {
	ctx := context.Background()
	_, err := s.passwordHistoryRepo.GeneratePasswordHistory(ctx, id, passwordHash)
	if err != nil {
		return fmt.Errorf("failed to add local user password history: %w", err)
	}
	return nil
}

func (s *localUserService) DeleteLocalUserPasswordHistory(id uuid.UUID) error {
	ctx := context.Background()
	err := s.passwordHistoryRepo.DeleteByUserID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete local user password history: %w", err)
	}
	return nil
}

func (s *localUserService) CreateLocalUserFlag(flag *models.LocalUserFlag) (uuid.UUID, error) {
	ctx := context.Background()
	createdFlagID, err := s.localUserFlagRepo.Insert(ctx, flag)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create local user flag: %w", err)
	}
	return createdFlagID, nil
}

func (s *localUserService) GetLocalUserFlag(id string) (models.LocalUserFlag, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.LocalUserFlag{}, fmt.Errorf("invalid local user flag ID: %w", err)
	}
	flag, err := s.localUserFlagRepo.Query(ctx, parsedID)
	if err != nil {
		return models.LocalUserFlag{}, fmt.Errorf("failed to get local user flag: %w", err)
	}
	return flag, nil
}

func (s *localUserService) UpdateLocalUserFlag(flag *models.LocalUserFlag) (models.LocalUserFlag, error) {
	ctx := context.Background()
	updatedFlag, err := s.localUserFlagRepo.Update(ctx, flag)
	if err != nil {
		return models.LocalUserFlag{}, fmt.Errorf("failed to update local user flag: %w", err)
	}
	return updatedFlag, nil
}

func (s *localUserService) DeleteLocalUserFlag(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid local user flag ID: %w", err)
	}
	err = s.localUserFlagRepo.Delete(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete local user flag: %w", err)
	}
	return nil
}

func (s *localUserService) GetLocalUserFlagByUserID(userID string) (models.LocalUserFlag, error) {
	ctx := context.Background()
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return models.LocalUserFlag{}, fmt.Errorf("invalid user ID: %w", err)
	}
	flag, err := s.localUserFlagRepo.QueryByUser(ctx, parsedUserID)
	if err != nil {
		return models.LocalUserFlag{}, fmt.Errorf("failed to get local user flag by user ID: %w", err)
	}
	return flag, nil
}

func (s *localUserService) CreateUserSecurityQuestion(question *models.UserSecurityQuestion)  error {
	ctx := context.Background()
	_, err := s.userSecurityQuestionRepo.Insert(ctx, question)
	if err != nil {
		return  fmt.Errorf("failed to create user security question: %w", err)
	}
	return  nil
}

func (s *localUserService) GetUserSecurityQuestion(id string) (models.UserSecurityQuestion, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.UserSecurityQuestion{}, fmt.Errorf("invalid user security question ID: %w", err)
	}
	question, err := s.userSecurityQuestionRepo.Query(ctx, parsedID)
	if err != nil {
		return models.UserSecurityQuestion{}, fmt.Errorf("failed to get user security question: %w", err)
	}
	return question, nil
}

func (s *localUserService) UpdateUserSecurityQuestion(question *models.UserSecurityQuestion) (models.UserSecurityQuestion, error) {
	ctx := context.Background()
	updatedQuestion, err := s.userSecurityQuestionRepo.Update(ctx, question)
	if err != nil {
		return models.UserSecurityQuestion{}, fmt.Errorf("failed to update user security question: %w", err)
	}
	return updatedQuestion, nil
}

func (s *localUserService) DeleteUserSecurityQuestion(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user security question ID: %w", err)
	}
	err = s.userSecurityQuestionRepo.Delete(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete user security question: %w", err)
	}
	return nil
}

func (s *localUserService) GetUserSecurityQuestionByUser(userID string) ([]models.UserSecurityQuestion, error) {
	ctx := context.Background()
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	questions, err := s.userSecurityQuestionRepo.QueryByUser(ctx, parsedUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user security questions by user ID: %w", err)
	}
	return questions, nil
}

func (s *localUserService) DeleteUserSecurityQuestionsByUser(userID string) error {
	ctx := context.Background()
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	err = s.userSecurityQuestionRepo.DeleteByUser(ctx, parsedUserID)
	if err != nil {
		return fmt.Errorf("failed to delete user security questions by user ID: %w", err)
	}
	return nil
}



