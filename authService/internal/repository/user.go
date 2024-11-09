package repository

import (
	"authService/config"
	"authService/db"
	"authService/internal/models"
	"context"

	"github.com/google/uuid"
)



type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
// Insert adds a new user to the database.
func (r *UserRepository) Insert(ctx context.Context, user models.User) <-chan ResultChan[models.User] {
	resultChan := make(chan ResultChan[models.User], 1)
	go func() {
		db.InitDB(config.GetDBConfig())
		defer close(resultChan)
		query := `INSERT INTO auth.users (username, auth_method, email) VALUES ($1, $2, $3) RETURNING id`
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
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) <-chan ResultChan[bool] {
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
func (r *UserRepository) Query(ctx context.Context, id uuid.UUID) <-chan ResultChan[models.User] {
	resultChan := make(chan ResultChan[models.User], 1)

	go func() {
		defer close(resultChan)
		query := `SELECT id, username, auth_method, email, status, email_verified, last_active_at FROM users WHERE id = $1`
		pool := db.GetPools().ReadPool // Use the ReadPool for select operations

		user := models.User{}
		err := pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.AuthMethod, &user.Email, &user.Status, &user.EmailVerified, &user.LastActiveAt)
		resultChan <- ResultChan[models.User]{Error: err, Data: user}
	}()

	return resultChan
}

// OAuthUserRepository defines methods for OAuth user data access.
type OAuthUserRepository struct{}

// NewOAuthUserRepository creates a new instance of OAuthUserRepository.
func NewOAuthUserRepository() *OAuthUserRepository {
	return &OAuthUserRepository{}
}

// Insert adds a new OAuth user to the database.
func (r *OAuthUserRepository) Insert(ctx context.Context, oauthUser models.OAuthUser) <-chan ResultChan[models.OAuthUser] {
	resultChan := make(chan ResultChan[models.OAuthUser], 1)

	go func() {
		defer close(resultChan)
		query := `INSERT INTO oauth_users (user_id, provider, provider_user_id, access_token, refresh_token, expires_at) 
		          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		pool := db.GetPools().CreatePool // Use the CreatePool for insert operations

		err := pool.QueryRow(ctx, query, oauthUser.UserID, oauthUser.Provider, oauthUser.ProviderUserID, oauthUser.AccessToken, oauthUser.RefreshToken, oauthUser.ExpiresAt).Scan(&oauthUser.ID)
		resultChan <- ResultChan[models.OAuthUser]{Error: err, Data: oauthUser}
	}()

	return resultChan
}

// Query retrieves an OAuth user by provider and provider user ID.
func (r *OAuthUserRepository) Query(ctx context.Context, provider models.OAuthProvider, providerUserID string) <-chan ResultChan[models.OAuthUser] {
	resultChan := make(chan ResultChan[models.OAuthUser], 1)

	go func() {
		defer close(resultChan)
		query := `SELECT id, user_id, provider, provider_user_id, access_token, refresh_token, expires_at 
		          FROM oauth_users WHERE provider = $1 AND provider_user_id = $2`
		pool := db.GetPools().ReadPool // Use the ReadPool for select operations

		oauthUser := models.OAuthUser{}
		err := pool.QueryRow(ctx, query, provider, providerUserID).Scan(&oauthUser.ID, &oauthUser.UserID, &oauthUser.Provider, &oauthUser.ProviderUserID, &oauthUser.AccessToken, &oauthUser.RefreshToken, &oauthUser.ExpiresAt)
		resultChan <- ResultChan[models.OAuthUser]{Error: err, Data: oauthUser}
	}()

	return resultChan
}

// LocalUserRepository defines methods for local user data access.
type LocalUserRepository struct{}

// NewLocalUserRepository creates a new instance of LocalUserRepository.
func NewLocalUserRepository() *LocalUserRepository {
	return &LocalUserRepository{}
}

// Insert adds a new local user to the database.
func (r *LocalUserRepository) Insert(ctx context.Context, localUser models.LocalUser) <-chan ResultChan[models.LocalUser] {
	resultChan := make(chan ResultChan[models.LocalUser], 1)

	go func() {
		defer close(resultChan)
		query := `INSERT INTO local_users (user_id, password, last_password_change, password_history, force_password_change, password_expires_at) 
		          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		pool := db.GetPools().CreatePool // Use the CreatePool for insert operations

		err := pool.QueryRow(ctx, query, localUser.UserID, localUser.Password, localUser.LastPasswordChange, localUser.PasswordHistory, localUser.ForcePasswordChange, localUser.PasswordExpiresAt).Scan(&localUser.ID)
		resultChan <- ResultChan[models.LocalUser]{Error: err, Data: localUser}
	}()

	return resultChan
}

// Query retrieves a local user by user ID.
func (r *LocalUserRepository) Query(ctx context.Context, userID uuid.UUID) <-chan ResultChan[models.LocalUser] {
	resultChan := make(chan ResultChan[models.LocalUser], 1)

	go func() {
		defer close(resultChan)
		query := `SELECT id, user_id, password, last_password_change, password_history, force_password_change, password_expires_at 
		          FROM local_users WHERE user_id = $1`
		pool := db.GetPools().ReadPool // Use the ReadPool for select operations

		localUser := models.LocalUser{}
		err := pool.QueryRow(ctx, query, userID).Scan(&localUser.ID, &localUser.UserID, &localUser.Password, &localUser.LastPasswordChange, &localUser.PasswordHistory, &localUser.ForcePasswordChange, &localUser.PasswordExpiresAt)
		resultChan <- ResultChan[models.LocalUser]{Error: err, Data: localUser}
	}()

	return resultChan
}
