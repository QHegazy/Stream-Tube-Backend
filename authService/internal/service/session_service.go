package services

import (
	"authService/internal/models"
	"authService/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SessionService interface {
	CreateSession(session *models.Session) (models.Session, error)
	GetSession(id string) (models.Session, error)
	UpdateSession(session *models.Session) (models.Session, error)
	DeleteSession(id string) error
	GetSessionByRefreshToken(refreshToken string) (models.Session, error)
	DeleteSessionByRefreshToken(refreshToken string) error
	UpdateSessionLastActivity(id string, lastActivity time.Time) error
	DeleteSessionsByUser(id string) error
	GetSessionsByUser(id string) ([]models.Session, error)
	BlacklistSession(id string, blacklistedAt time.Time) error
	UnblacklistSession(id string) error
	CountActiveSessionsByUser(id string) (int, error)
	DeleteExpiredSessions() error
}

type sessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService() SessionService {
	return &sessionService{
		sessionRepo: *repository.NewSessionRepository(),
	}
}

func (s *sessionService) CreateSession(session *models.Session) (models.Session, error) {
	ctx := context.Background()
	createdSession, err := s.sessionRepo.Insert(ctx, session)
	if err != nil {
		return models.Session{}, fmt.Errorf("failed to create session: %w", err)
	}
	return createdSession, nil
}

func (s *sessionService) GetSession(id string) (models.Session, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Session{}, fmt.Errorf("invalid session ID: %w", err)
	}
	session, err := s.sessionRepo.Query(ctx, parsedID)
	if err != nil {
		return models.Session{}, fmt.Errorf("failed to get session: %w", err)
	}
	return session, nil
}

func (s *sessionService) UpdateSession(session *models.Session) (models.Session, error) {
	ctx := context.Background()
	updatedSession, err := s.sessionRepo.Update(ctx, session)
	if err != nil {
		return models.Session{}, fmt.Errorf("failed to update session: %w", err)
	}
	return updatedSession, nil
}

func (s *sessionService) DeleteSession(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}
	err = s.sessionRepo.Delete(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

func (s *sessionService) GetSessionByRefreshToken(refreshToken string) (models.Session, error) {
	ctx := context.Background()
	session, err := s.sessionRepo.QueryByRefreshToken(ctx, refreshToken)
	if err != nil {
		return models.Session{}, fmt.Errorf("failed to get session by refresh token: %w", err)
	}
	return session, nil
}

func (s *sessionService) DeleteSessionByRefreshToken(refreshToken string) error {
	ctx := context.Background()
	err := s.sessionRepo.DeleteByRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to delete session by refresh token: %w", err)
	}
	return nil
}

func (s *sessionService) UpdateSessionLastActivity(id string, lastActivity time.Time) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}
	err = s.sessionRepo.UpdateLastActivity(ctx, parsedID, lastActivity)
	if err != nil {
		return fmt.Errorf("failed to update session last activity: %w", err)
	}
	return nil
}

func (s *sessionService) DeleteSessionsByUser(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}
	err = s.sessionRepo.DeleteByUser(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to delete sessions by user ID: %w", err)
	}
	return nil
}

func (s *sessionService) GetSessionsByUser(id string) ([]models.Session, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	sessions, err := s.sessionRepo.QueryByUser(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions by user ID: %w", err)
	}
	return sessions, nil
}

func (s *sessionService) BlacklistSession(id string, blacklistedAt time.Time) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}
	err = s.sessionRepo.BlacklistSession(ctx, parsedID, blacklistedAt)
	if err != nil {
		return fmt.Errorf("failed to blacklist session: %w", err)
	}
	return nil
}

func (s *sessionService) UnblacklistSession(id string) error {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid session ID: %w", err)
	}
	err = s.sessionRepo.UnblacklistSession(ctx, parsedID)
	if err != nil {
		return fmt.Errorf("failed to unblacklist session: %w", err)
	}
	return nil
}

func (s *sessionService) CountActiveSessionsByUser(id string) (int, error) {
	ctx := context.Background()
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}
	count, err := s.sessionRepo.CountActiveSessionsByUser(ctx, parsedID)
	if err != nil {
		return 0, fmt.Errorf("failed to count active sessions by user ID: %w", err)
	}
	return count, nil
}

func (s *sessionService) DeleteExpiredSessions() error {
	ctx := context.Background()
	err := s.sessionRepo.DeleteExpiredSessions(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}
	return nil
}
