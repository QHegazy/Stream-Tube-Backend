package models

import (
	"time"

	"github.com/google/uuid"
)




type User struct {
    ID             uuid.UUID  `json:"id" db:"id" default:"uuid_generate_v4()"`
    Username       string     `json:"username" db:"username" validate:"required,min=3,max=50,alphanumunicode"`
    AuthMethod     AuthMethod `json:"auth_method" db:"auth_method"`
    Email          string     `json:"email" db:"email" validate:"required,email"`
    Status         UserStatus `json:"status" db:"status" default:"pending_verification"` 
    EmailVerified  *time.Time `json:"email_verified,omitempty" db:"email_verified"`
    LastActiveAt   *time.Time `json:"last_active_at,omitempty" db:"last_active_at"`
}




// OAuthUsers Table
type OAuthUser struct {
    ID             uuid.UUID     `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID         uuid.UUID     `json:"user_id" db:"user_id"`
    Provider       OAuthProvider `json:"provider" db:"provider"`
    ProviderUserID string        `json:"provider_user_id" db:"provider_user_id"`
    AccessToken    string        `json:"access_token" db:"access_token"`
    RefreshToken   *string       `json:"refresh_token,omitempty" db:"refresh_token"`
    ExpiresAt      *time.Time    `json:"expires_at,omitempty" db:"expires_at"`
}

// LocalUsers Table
type LocalUser struct {
    ID                uuid.UUID   `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID            uuid.UUID   `json:"user_id" db:"user_id"`
    Password          string      `json:"password" db:"password"`
    LastPasswordChange *time.Time `json:"last_password_change,omitempty" db:"last_password_change"`
    PasswordHistory   []string    `json:"password_history" db:"password_history"`
    ForcePasswordChange bool      `json:"force_password_change" db:"force_password_change"`
    PasswordExpiresAt *time.Time  `json:"password_expires_at,omitempty" db:"password_expires_at"`
}
