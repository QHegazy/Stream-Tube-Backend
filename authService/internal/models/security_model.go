package models

import (
	"time"

	"github.com/google/uuid"
)


type MFASettings struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	MFAType   MFAType   `json:"mfa_type" db:"mfa_type"`
	LastMFAAt *time.Time `json:"last_mfa_at" db:"last_mfa_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type MFATOTP struct {
	ID            uuid.UUID `json:"id" db:"id"`                     
	MFASettingsID uuid.UUID `json:"mfa_settings_id" db:"mfa_settings_id"`
	Secret        string    `json:"secret" db:"secret"`              
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type BackupCode struct {
	ID         uuid.UUID `json:"id" db:"id"`
	MFATOTPID  uuid.UUID `json:"mfa_totp_id" db:"mfa_totpid"`
	Code       string    `json:"code" db:"code"`                   
	UsedCount  int       `json:"used_count" db:"used_count"`       
	CreatedAt  time.Time `json:"created_at" db:"created_at"` 
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"` 
}
type AccountSecurityStatus struct {
	StatusID                uuid.UUID      `json:"status_id" db:"status_id"`
	UserID                  uuid.UUID      `json:"user_id" db:"user_id"`
	RiskLevel               string         `json:"risk_level" db:"risk_level"` 
	LockReason              *string        `json:"lock_reason" db:"lock_reason"`
	LockedUntil             *time.Time     `json:"locked_until" db:"locked_until"`
	FailedLoginAttempts     int            `json:"failed_login_attempts" db:"failed_login_attempts"`
	FailedMFAAttempts       int            `json:"failed_mfa_attempts" db:"failed_mfa_attempts"`
	SuspiciousActivityCount int            `json:"suspicious_activity_count" db:"suspicious_activity_count"`
	FailedLoginResetAt     *time.Time     `json:"failed_login_reset_at" db:"failed_login_reset_at"`
	LastLoginAt            *time.Time     `json:"last_login_at" db:"last_login_at"`
	LastLoginIP            *string        `json:"last_login_ip" db:"last_login_ip"`
	KnownDevices           []interface{}  `json:"known_devices" db:"known_devices"` 
	TrustedLocations       []interface{}  `json:"trusted_locations" db:"trusted_locations"` 
	CreatedAt              time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at" db:"updated_at"`
}

type UserAccessLog struct {
	LogID             uuid.UUID   `json:"log_id" db:"log_id"`
	UserID            uuid.UUID   `json:"user_id" db:"user_id"`
	SessionID         *uuid.UUID  `json:"session_id" db:"session_id"`
	AuthMethod        string      `json:"auth_method" db:"auth_method"`
	AuthProvider      *string     `json:"auth_provider" db:"auth_provider"`
	IPAddress         string      `json:"ip_address" db:"ip_address"`
	UserAgent         string      `json:"user_agent" db:"user_agent"`
	DeviceInfo        string      `json:"device_info" db:"device_info"`
	LocationInfo      string      `json:"location_info" db:"location_info"`
	LoginSuccessful   bool        `json:"login_successful" db:"login_successful"`
	MFAUsed           bool        `json:"mfa_used" db:"mfa_used"`
	MFAType           *string     `json:"mfa_type" db:"mfa_type"`
	RiskScore         *int        `json:"risk_score" db:"risk_score"`
	Metadata          string      `json:"metadata" db:"metadata"`
	FailureReason     *string     `json:"failure_reason" db:"failure_reason"`
	LoginTimestamp    time.Time   `json:"login_timestamp" db:"login_timestamp"`
}

type UserPreferences struct {
	ID                        uuid.UUID      `json:"id" db:"id"`
	UserID                    uuid.UUID      `json:"user_id" db:"user_id"`
	EmailNotifications        bool           `json:"email_notifications" db:"email_notifications"`
	SMSNotifications          bool           `json:"sms_notifications" db:"sms_notifications"`
	PushNotifications         bool           `json:"push_notifications" db:"push_notifications"`
	TwoFactorAuthEnabled      bool           `json:"two_factor_auth_enabled" db:"two_factor_auth_enabled"`
	PreferredCommunicationMethod string      `json:"preferred_communication_method" db:"preferred_communication_method"`
	MarketingPreferences      string         `json:"marketing_preferences" db:"marketing_preferences"` 
	PrivacySettings           string         `json:"privacy_settings" db:"privacy_settings"`         
	ThemePreference           string         `json:"theme_preference" db:"theme_preference"`

}


type NotificationTemplate struct {
	ID        uuid.UUID       `json:"id" db:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Type      NotificationType `json:"type" db:"type,notnull"`
	Name      string          `json:"name" db:"name,notnull"`
	Subject   *string         `json:"subject" db:"subject"`
	Content   string          `json:"content" db:"content,notnull"`
}

type NotificationHistory struct {
	ID               uuid.UUID        `json:"id" db:"id,pk,type:uuid,default:uuid_generate_v4()"`
	UserID           uuid.UUID        `json:"user_id" db:"user_id,type:uuid,notnull"`
	TemplateID       uuid.UUID        `json:"template_id" db:"template_id,type:uuid,notnull"`
	NotificationType NotificationType `json:"notification_type" db:"notification_type,type:notification_type,notnull"`
	SentAt           time.Time        `json:"sent_at" db:"sent_at,type:timestamptz,default:now()"`
	DeliveredAt      *time.Time       `json:"delivered_at" db:"delivered_at,type:timestamptz"`
	ReadAt           *time.Time       `json:"read_at" db:"read_at,type:timestamptz"`
}


type RateLimit struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	UserID        uuid.UUID     `json:"user_id" db:"user_id"`
	IpAddress     string        `json:"ip_address" db:"ip_address"`
	Endpoint      string        `json:"endpoint" db:"endpoint"`
	RequestCount  int           `json:"request_count" db:"request_count"`
	WindowStart   time.Time     `json:"window_start" db:"window_start"`
	WindowDuration time.Duration`json:"window_duration" db:"window_duration"`
	MaxRequests   int           `json:"max_requests" db:"max_requests"`
}

type Device struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	OS         string    `json:"os" db:"os"`
	Browser    string    `json:"browser" db:"browser"`
	DeviceType string    `json:"device_type" db:"device_type"`
	IPAddress  string    `json:"ip_address" db:"ip_address"`
	UserAgent  string   `json:"user_agent" db:"user_agent,omitempty"`
	LastUsedAt time.Time `json:"last_used_at" db:"last_used_at,omitempty"`
	Region     string   `json:"region" db:"region,omitempty"`
}