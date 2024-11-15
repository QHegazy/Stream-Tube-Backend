package repository

import (
	"authService/config"
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)



type UserRepository struct{
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
// Insert inserts a new user into the database.
func (r *UserRepository) Insert(ctx context.Context, user *models.User, resultChan *chan ResultChan[models.User]) {
    go func() {
        
        db.InitDB(config.GetDBConfig())
        query := `INSERT INTO auth.users (username, auth_method, email, status, email_verified, last_active_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
        pool := db.GetPools().CreatePool

        err := pool.QueryRow(ctx, query, user.Username, user.AuthMethod, user.Email, user.Status, user.EmailVerified, user.LastActiveAt).Scan(&user.ID)
        if err != nil {
            *resultChan <- ResultChan[models.User]{
                Error: fmt.Errorf("failed to insert user: %w", err),
                Data:  models.User{}, // Empty user struct instead of nil
            }
            return
        }

        *resultChan <- ResultChan[models.User]{
            Error: nil,
            Data:  *user,
        }
    }()
}
// Update updates an existing user in the database.
func (r *UserRepository) Update(ctx context.Context, user *models.User, resultChan *chan ResultChan[models.User]) {
    go func() {
        defer close(*resultChan)
        
        query := `
            UPDATE auth.users 
            SET username = $1, 
                email = $2, 
                auth_method = $3,
                updated_at = $4
            WHERE id = $5
            RETURNING id, username, auth_method, email, status, email_verified, last_active_at`
            
        pool := db.GetPools().UpdatePool

        updatedUser := models.User{}
        err := pool.QueryRow(
            ctx,
            query,
            user.Username,
            user.Email,
            user.AuthMethod,
            time.Now(),
            user.ID,
        ).Scan(
            &updatedUser.ID,
            &updatedUser.Username,
            &updatedUser.AuthMethod,
            &updatedUser.Email,
            &updatedUser.Status,
            &updatedUser.EmailVerified,
            &updatedUser.LastActiveAt,
        )

        if err != nil {
            *resultChan <- ResultChan[models.User]{
                Error: fmt.Errorf("failed to update user: %w", err),
                Data:  models.User{},
            }
            return
        }

        *resultChan <- ResultChan[models.User]{
            Error: nil,
            Data:  updatedUser,
        }
    }()
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID, resultChan *chan ResultChan[models.User]) {
    go func() {
        defer close(*resultChan)
        
        query := `
            DELETE FROM auth.users 
            WHERE id = $1 
            RETURNING id, username, auth_method, email, status, email_verified, last_active_at`
            
        pool := db.GetPools().DeletePool

        deletedUser := models.User{}
        err := pool.QueryRow(ctx, query, id).Scan(
            &deletedUser.ID,
            &deletedUser.Username,
            &deletedUser.AuthMethod,
            &deletedUser.Email,
            &deletedUser.Status,
            &deletedUser.EmailVerified,
            &deletedUser.LastActiveAt,
        )

        if err != nil {
            *resultChan <- ResultChan[models.User]{
                Error: fmt.Errorf("failed to delete user: %w", err),
                Data:  models.User{},
            }
            return
        }

        *resultChan <- ResultChan[models.User]{
            Error: nil,
            Data:  deletedUser,
        }
    }()
}

func (r *UserRepository) Query(ctx context.Context, id uuid.UUID, resultChan *chan ResultChan[models.User]) {
    go func() {
        defer close(*resultChan)
        
        query := `
            SELECT 
                id, 
                username, 
                auth_method, 
                email, 
                status, 
                email_verified, 
                last_active_at,
                created_at,
                updated_at
            FROM auth.users 
            WHERE id = $1`
            
        pool := db.GetPools().ReadPool

        user := models.User{}
        err := pool.QueryRow(ctx, query, id).Scan(
            &user.ID,
            &user.Username,
            &user.AuthMethod,
            &user.Email,
            &user.Status,
            &user.EmailVerified,
            &user.LastActiveAt,
        )

        if err != nil {
            *resultChan <- ResultChan[models.User]{
                Error: fmt.Errorf("failed to query user: %w", err),
                Data:  models.User{},
            }
            return
        }

        *resultChan <- ResultChan[models.User]{
            Error: nil,
            Data:  user,
        }
    }()
}


func (r *UserRepository) Exists(ctx context.Context, id uuid.UUID) bool {
    query := `SELECT EXISTS(SELECT 1 FROM auth.users WHERE id = $1)`
    pool := db.GetPools().ReadPool
    
    var exists bool
    err := pool.QueryRow(ctx, query, id).Scan(&exists)
    if err != nil {
        return false
    }
    return exists
}