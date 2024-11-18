package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UserRepository implements Operations for models.User.
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Insert inserts a new user into the database.
func (r *UserRepository) Insert(ctx context.Context, user models.User) (*models.User, error) {
	query := `INSERT INTO auth.users (username, auth_method, email, status, email_verified, last_active_at) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	pool := db.GetPools().CreatePool

	err := pool.QueryRow(ctx, query, user.Username, user.AuthMethod, user.Email, user.Status, user.EmailVerified, user.LastActiveAt).
		Scan(&user.ID)
	if err != nil {
		return &models.User{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return &user, nil
}

// Update updates an existing user in the database.
func (r *UserRepository) Update(ctx context.Context, user models.User) (models.User, error) {
	query := `UPDATE auth.users 
              SET username = $1, email = $2, auth_method = $3, updated_at = $4 
              WHERE id = $5 
              RETURNING id, username, auth_method, email, status, email_verified, last_active_at`
	pool := db.GetPools().UpdatePool

	var updatedUser models.User
	err := pool.QueryRow(ctx, query, user.Username, user.Email, user.AuthMethod, time.Now(), user.ID).
		Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.AuthMethod, &updatedUser.Email, &updatedUser.Status, &updatedUser.EmailVerified, &updatedUser.LastActiveAt)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

// Delete removes a user from the database by ID.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM auth.users WHERE id = $1`
	pool := db.GetPools().DeletePool

	_, err := pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// Query retrieves a user from the database by ID.
func (r *UserRepository) Query(ctx context.Context, id uuid.UUID) (models.User, error) {
	query := `SELECT id, username, auth_method, email, status, email_verified, last_active_at, created_at, updated_at 
              FROM auth.users 
              WHERE id = $1`
	pool := db.GetPools().ReadPool

	var user models.User
	err := pool.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.AuthMethod, &user.Email, &user.Status, &user.EmailVerified, &user.LastActiveAt)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to query user: %w", err)
	}

	return user, nil
}
