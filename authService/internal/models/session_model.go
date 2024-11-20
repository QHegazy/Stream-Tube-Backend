package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID            uuid.UUID  `json:"id"`             
	UserID        uuid.UUID  `json:"user_id"`        
    DeviceID      uuid.UUID  `json:"device_id"`      
	RefreshToken  string     `json:"refresh_token"`  
	LastActivity  time.Time  `json:"last_activity"`  
	BlacklistedAt time.Time `json:"blacklisted_at"` 
	ExpiresAt     time.Time  `json:"expires_at"`     
}
