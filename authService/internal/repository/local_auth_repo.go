package repository

import (
	"authService/config"
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type LocalAuthRepository struct {
}

func NewLocalAuthRepository() *LocalAuthRepository {
	return &LocalAuthRepository{}
}

func (r *LocalAuthRepository) Insert(ctx context.Context, localUser *models.LocalUser, resultChan *chan ResultChan[models.LocalUser]) {
	go func() {
		db.InitDB(config.GetDBConfig())
		query := `INSERT INTO auth.local_users (user_id, password, last_password_change, password_history, force_password_change, password_expires_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		pool := db.GetPools().CreatePool

		err := pool.QueryRow(ctx, query, localUser.UserID, localUser.Password, localUser.LastPasswordChange, localUser.PasswordHistory, localUser.ForcePasswordChange, localUser.PasswordExpiresAt).Scan(&localUser.ID)
		if err != nil {
			*resultChan <- ResultChan[models.LocalUser]{
				Error: fmt.Errorf("failed to insert local user: %w", err),
				Data:  models.LocalUser{},
			}
			return
		}

		*resultChan <- ResultChan[models.LocalUser]{
			Error: nil,
			Data:  *localUser,
		}
	}()
}

func (r *LocalAuthRepository) Query(ctx context.Context, id uuid.UUID, resultChan *chan ResultChan[models.LocalUser]) {
	go func() {
		db.InitDB(config.GetDBConfig())
		query := `SELECT id, user_id, password, last_password_change, password_history, force_password_change, password_expires_at FROM auth.local_users WHERE id = $1`
		pool := db.GetPools().ReadPool

		localUser := models.LocalUser{}
		err := pool.QueryRow(ctx, query, id).Scan(
			&localUser.ID,
			&localUser.UserID,
			&localUser.Password,
			&localUser.LastPasswordChange,
			&localUser.PasswordHistory,
			&localUser.ForcePasswordChange,
			&localUser.PasswordExpiresAt,
		)
		if err != nil {
			*resultChan <- ResultChan[models.LocalUser]{
				Error: fmt.Errorf("failed to query local user: %w", err),
				Data:  models.LocalUser{},
			}
			return
		}

		*resultChan <- ResultChan[models.LocalUser]{
			Error: nil,
			Data:  localUser,
		}
	}()
}

func (r *LocalAuthRepository) Update(ctx context.Context, localUser *models.LocalUser, resultChan *chan ResultChan[models.LocalUser]) {
	go func() {
		db.InitDB(config.GetDBConfig())
		query := `UPDATE auth.local_users SET user_id = $1, password = $2, last_password_change = $3, password_history = $4, force_password_change = $5, password_expires_at = $6 WHERE id = $7 RETURNING id`
		pool := db.GetPools().UpdatePool

		updatedLocalUser := models.LocalUser{}
		err := pool.QueryRow(ctx, query, localUser.UserID, localUser.Password, localUser.LastPasswordChange, localUser.PasswordHistory, localUser.ForcePasswordChange, localUser.PasswordExpiresAt, localUser.ID).Scan(&updatedLocalUser.ID)
		if err != nil {
			*resultChan <- ResultChan[models.LocalUser]{
				Error: fmt.Errorf("failed to update local user: %w", err),
				Data:  models.LocalUser{},
			}
			return
		}

		*resultChan <- ResultChan[models.LocalUser]{
			Error: nil,
			Data:  updatedLocalUser,
		}
	}()
}

func (r *LocalAuthRepository) Delete(ctx context.Context, id uuid.UUID, resultChan *chan ResultChan[models.LocalUser]) {
	go func() {
		db.InitDB(config.GetDBConfig())
		query := `DELETE FROM auth.local_users WHERE id = $1 RETURNING id`
		pool := db.GetPools().DeletePool

		deletedLocalUser := models.LocalUser{}
		err := pool.QueryRow(ctx, query, id).Scan(&deletedLocalUser.ID)
		if err != nil {
			*resultChan <- ResultChan[models.LocalUser]{
				Error: fmt.Errorf("failed to delete local user: %w", err),
				Data:  models.LocalUser{},
			}
			return
		}

		*resultChan <- ResultChan[models.LocalUser]{
			Error: nil,
			Data:  deletedLocalUser,
		}
	}()
}
func (r *LocalAuthRepository) QueryByEmail(ctx context.Context, email string, resultChan *chan ResultChan[models.LocalUser]) {
	go func() {
		db.InitDB(config.GetDBConfig())
		query := `SELECT lu.id, lu.user_id, lu.password, lu.last_password_change, lu.password_history, lu.force_password_change, lu.password_expires_at
		FROM auth.local_users lu
		JOIN auth.users u ON lu.user_id = u.id
		WHERE u.email = $1`
		pool := db.GetPools().ReadPool

		localUser := models.LocalUser{}
		err := pool.QueryRow(ctx, query, email).Scan(
			&localUser.ID,
			&localUser.UserID,
			&localUser.Password,
			&localUser.LastPasswordChange,
			&localUser.PasswordHistory,
			&localUser.ForcePasswordChange,
			&localUser.PasswordExpiresAt,
		)
		if err != nil {
			*resultChan <- ResultChan[models.LocalUser]{
				Error: fmt.Errorf("failed to query local user by email: %w", err),
				Data:  models.LocalUser{},
			}
			return
		}

		*resultChan <- ResultChan[models.LocalUser]{
			Error: nil,
			Data:  localUser,
		}
	}()
}
