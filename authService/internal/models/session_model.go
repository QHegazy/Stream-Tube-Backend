package models

import (
	"time"

	"github.com/google/uuid"
)

// Session Table
type Session struct {
    ID               uuid.UUID  `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID           uuid.UUID  `json:"user_id" db:"user_id"`
    AccessToken      string     `json:"access_token" db:"access_token"`
    RefreshToken     *string    `json:"refresh_token,omitempty" db:"refresh_token"`
    ExpiresAt        time.Time  `json:"expires_at" db:"expires_at"`
    IPAddress        string     `json:"ip_address" db:"ip_address"`
    UserAgent        string     `json:"user_agent" db:"user_agent"`
    DeviceInfo       string     `json:"device_info,omitempty" db:"device_info"`
    LoginTime        time.Time  `json:"login_time" db:"login_time"`
    LogoutTime       *time.Time `json:"logout_time,omitempty" db:"logout_time"`
}

