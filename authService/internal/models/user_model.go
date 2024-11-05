package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthMethod string

const (
    AuthMethodLocal AuthMethod = "local"
    AuthMethodOAuth AuthMethod = "oauth"
    AuthMethodMFA   AuthMethod = "mfa"
)

type User struct {
    ID             uuid.UUID  `json:"id" db:"id" default:"uuid_generate_v4()"`
    Username       string     `json:"username" db:"username" validate:"required,min=3,max=50,alphanumunicode"`
    AuthMethod     AuthMethod `json:"auth_method" db:"auth_method"`
    Email          string     `json:"email" db:"email" validate:"required,email"`
    Status         string     `json:"status" db:"status" default:"pending_verification"` 
    EmailVerified  *time.Time `json:"email_verified,omitempty" db:"email_verified"`
    LastActiveAt   *time.Time `json:"last_active_at,omitempty" db:"last_active_at"`
}

func (a AuthMethod) IsValid() bool {
    switch a {
    case AuthMethodLocal, AuthMethodOAuth, AuthMethodMFA:
        return true
    }
    return false
}
