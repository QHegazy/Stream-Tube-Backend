package models

import (
	"time"

	"github.com/google/uuid"
)




type User struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	Username       string         `json:"username" db:"username"`
	AuthMethod     AuthMethod     `json:"auth_method" db:"auth_method"`
	Email          string         `json:"email" db:"email"`
	Status         UserStatus     `json:"status" db:"status"`
	EmailVerified  time.Time      `json:"email_verified" db:"email_verified"`
	LastActiveAt   time.Time      `json:"last_active_at" db:"last_active_at"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time     `json:"deleted_at" db:"deleted_at"`
}




// OAuthUsers Table
type OAuthUser struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	UserID         uuid.UUID      `json:"user_id" db:"user_id"`
	Provider       OAuthProvider  `json:"provider" db:"provider"`
	ProviderUserID string         `json:"provider_user_id" db:"provider_user_id"`
	AccessToken    string         `json:"access_token" db:"access_token"`
	RefreshToken   string         `json:"refresh_token" db:"refresh_token"`
	ExpiresAt      time.Time      `json:"expires_at" db:"expires_at"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time     `json:"deleted_at" db:"deleted_at"`
}

// LocalUsers Table
type LocalUser struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	UserID       uuid.UUID      `json:"user_id" db:"user_id"`
	PasswordHash string         `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at" db:"updated_at"`
}
type PasswordHistory struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	UserID       uuid.UUID      `json:"user_id" db:"user_id"`
	PasswordHash string         `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
}


type LocalUserFlag struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	LocalUserID        uuid.UUID `json:"local_user_id" db:"local_user_id"`
	ForcePasswordChange bool      `json:"force_password_change" db:"force_password_change"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type UserSecurityQuestion struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Question  string    `json:"question" db:"question"`
	Answer    string    `json:"answer" db:"answer"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
