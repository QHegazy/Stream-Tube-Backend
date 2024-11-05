package repository

import (
	"authService/db"
	"authService/internal/models"
	"context"
)

// ResultChan is a generic struct to encapsulate the result of operations.

// UserRepository defines methods for user data access.
type UserRepository struct{}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Insert adds a new user to the database.
func (r *UserRepository) Insert(ctx context.Context, user models.User) <-chan ResultChan[models.User] {
	resultChan := make(chan ResultChan[models.User], 1)

	go func() {
		defer close(resultChan)
		query := `INSERT INTO users (username, auth_method, email) VALUES ($1, $2, $3) RETURNING id`
		pool := db.GetPools().CreatePool // Use the CreatePool for insert operations

		err := pool.QueryRow(ctx, query, user.Username, user.AuthMethod, user.Email).Scan(&user.ID)
		resultChan <- ResultChan[models.User]{Error: err, Data: user}
	}()

	return resultChan
}

// Update updates an existing user in the database.
func (r *UserRepository) Update(ctx context.Context, user models.User) <-chan ResultChan[models.User] {
	resultChan := make(chan ResultChan[models.User], 1)

	go func() {
		defer close(resultChan)
		query := `UPDATE users SET username = $1, email = $2, auth_method = $3 WHERE id = $4`
		pool := db.GetPools().UpdatePool // Use the UpdatePool for update operations

		_, err := pool.Exec(ctx, query, user.Username, user.Email, user.AuthMethod, user.ID)
		resultChan <- ResultChan[models.User]{Error: err, Data: user}
	}()

	return resultChan
}

// Delete removes a user from the database by ID.
func (r *UserRepository) Delete(ctx context.Context, id int64) <-chan ResultChan[bool] {
	resultChan := make(chan ResultChan[bool], 1)

	go func() {
		defer close(resultChan)
		query := `DELETE FROM users WHERE id = $1`
		pool := db.GetPools().DeletePool // Use the DeletePool for delete operations

		_, err := pool.Exec(ctx, query, id)
		resultChan <- ResultChan[bool]{Error: err, Data: err == nil}
	}()

	return resultChan
}

// Query retrieves a user by ID.
func (r *UserRepository) Query(ctx context.Context, id int64) <-chan ResultChan[models.User] {
	resultChan := make(chan ResultChan[models.User], 1)

	go func() {
		defer close(resultChan)
		query := `SELECT id, username, auth_method, email FROM users WHERE id = $1`
		pool := db.GetPools().ReadPool // Use the ReadPool for select operations

		user := models.User{}
		err := pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.AuthMethod, &user.Email)
		resultChan <- ResultChan[models.User]{Error: err, Data: user}
	}()

	return resultChan
}
