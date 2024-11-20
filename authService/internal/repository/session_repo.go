package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type SessionRepository struct{}

func NewSessionRepository() *SessionRepository {
	return &SessionRepository{}
}

func (r *SessionRepository) Insert(ctx context.Context, session *models.Session) (models.Session, error) {
	query := `
		INSERT INTO auth.sessions (user_id, device_id, refresh_token, last_activity, blacklisted_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, session.UserID, session.DeviceID, session.RefreshToken, session.LastActivity, session.BlacklistedAt, session.ExpiresAt).Scan(&session.ID)
	if err != nil {
		return models.Session{}, fmt.Errorf("failed to insert session: %w", err)
	}

	return *session, nil
}

func (r *SessionRepository) Query(ctx context.Context, id uuid.UUID) (models.Session, error) {
	query := `
		SELECT id, user_id, device_id, refresh_token, last_activity, blacklisted_at, expires_at
		FROM auth.sessions WHERE id = $1`
	pool := db.GetPools().ReadPool

	session := models.Session{}
	err := pool.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.DeviceID,
		&session.RefreshToken,
		&session.LastActivity,
		&session.BlacklistedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		return models.Session{}, fmt.Errorf("failed to query session: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) Update(ctx context.Context, session *models.Session) (models.Session, error) {
	query := `
		UPDATE auth.sessions
		SET user_id = $1, device_id = $2, refresh_token = $3, last_activity = $4, blacklisted_at = $5, expires_at = $6
		WHERE id = $7 RETURNING id`
	pool := db.GetPools().UpdatePool

	updatedSession := models.Session{}
	err := pool.QueryRow(ctx, query, session.UserID, session.DeviceID, session.RefreshToken, session.LastActivity, session.BlacklistedAt, session.ExpiresAt, session.ID).Scan(&updatedSession.ID)
	if err != nil {
		return models.Session{}, fmt.Errorf("failed to update session: %w", err)
	}

	return updatedSession, nil
}

func (r *SessionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.sessions WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (r *SessionRepository) QueryByRefreshToken(ctx context.Context, refreshToken string) (models.Session, error) {
	query := `SELECT id, user_id, device_id, refresh_token, last_activity, blacklisted_at, expires_at FROM auth.sessions WHERE refresh_token = $1`
	pool := db.GetPools().ReadPool

	session := models.Session{}
	err := pool.QueryRow(ctx, query, refreshToken).Scan(
		&session.ID,
		&session.UserID,
		&session.DeviceID,
		&session.RefreshToken,
		&session.LastActivity,
		&session.BlacklistedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		return models.Session{}, fmt.Errorf("failed to query session by refresh token: %w", err)
	}

	return session, nil
}

func (r *SessionRepository) DeleteByRefreshToken(ctx context.Context, refreshToken string) error {
	query := `DELETE FROM auth.sessions WHERE refresh_token = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to delete session by refresh token: %w", err)
	}

	return nil
}
func (r *SessionRepository) UpdateLastActivity(ctx context.Context, id uuid.UUID, lastActivity time.Time) error {
	query := `UPDATE auth.sessions SET last_activity = $1 WHERE id = $2`
	pool := db.GetPools().UpdatePool

	_, err := pool.Exec(ctx, query, lastActivity, id)
	if err != nil {
		return fmt.Errorf("failed to update session last activity: %w", err)
	}

	return nil
}

func (r *SessionRepository) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM auth.sessions WHERE user_id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete sessions by user ID: %w", err)
	}

	return nil
}
func (r *SessionRepository) QueryByUser(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	query := `SELECT id, user_id, device_id, refresh_token, last_activity, blacklisted_at, expires_at FROM auth.sessions WHERE user_id = $1`
	pool := db.GetPools().ReadPool

	rows, err := pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sessions by user ID: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.DeviceID,
			&session.RefreshToken,
			&session.LastActivity,
			&session.BlacklistedAt,
			&session.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session row: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over session rows: %w", err)
	}

	return sessions, nil
}

func (r *SessionRepository) BlacklistSession(ctx context.Context, id uuid.UUID, blacklistedAt time.Time) error {
	query := `UPDATE auth.sessions SET blacklisted_at = $1 WHERE id = $2`
	pool := db.GetPools().UpdatePool

	_, err := pool.Exec(ctx, query, blacklistedAt, id)
	if err != nil {
		return fmt.Errorf("failed to blacklist session: %w", err)
	}

	return nil
}
func (r *SessionRepository) UnblacklistSession(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE auth.sessions SET blacklisted_at = NULL WHERE id = $1`
	pool := db.GetPools().UpdatePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to unblacklist session: %w", err)
	}

	return nil
}
func (r *SessionRepository) CountActiveSessionsByUser(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM auth.sessions WHERE user_id = $1 AND blacklisted_at IS NULL`
	pool := db.GetPools().ReadPool

	var count int
	err := pool.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active sessions by user ID: %w", err)
	}

	return count, nil
}
func (r *SessionRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM auth.sessions WHERE expires_at < NOW()`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}
