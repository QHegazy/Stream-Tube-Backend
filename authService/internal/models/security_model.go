package models

import (
	"time"

	"github.com/google/uuid"
)


type MFASecurity struct {
    ID          uuid.UUID  `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID      uuid.UUID  `json:"user_id" db:"user_id"`
    MFAType     MFAType    `json:"mfa_type" db:"mfa_type"`
    Secret      string     `json:"secret" db:"secret"`
    Enabled     bool       `json:"enabled" db:"enabled"`
}

// Account Security Status Table
type AccountSecurityStatus struct {
    ID               uuid.UUID `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID           uuid.UUID `json:"user_id" db:"user_id"`
    PasswordStatus   string    `json:"password_status" db:"password_status"`
    MFASecurityLevel string    `json:"mfa_security_level" db:"mfa_security_level"`
    LastRiskCheck    *time.Time `json:"last_risk_check,omitempty" db:"last_risk_check"`
    RiskLevel        RiskLevel  `json:"risk_level" db:"risk_level"`
}

// User Access Log Table
type UserAccessLog struct {
    ID           uuid.UUID `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID       uuid.UUID `json:"user_id" db:"user_id"`
    IPAddress    string    `json:"ip_address" db:"ip_address"`
    Location     string    `json:"location,omitempty" db:"location"`
    DeviceInfo   string    `json:"device_info" db:"device_info"`
    Action       string    `json:"action" db:"action"`
}

// User Preferences Table
type UserPreference struct {
    ID                  uuid.UUID `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID              uuid.UUID `json:"user_id" db:"user_id"`
    Language            string    `json:"language" db:"language"`
    TimeZone            string    `json:"time_zone" db:"time_zone"`
    NotificationEnabled bool      `json:"notification_enabled" db:"notification_enabled"`

}

// Notification Template Table
type NotificationTemplate struct {
    ID            uuid.UUID        `json:"id" db:"id" default:"uuid_generate_v4()"`
    Type          NotificationType `json:"type" db:"type"`
    TemplateName  string           `json:"template_name" db:"template_name"`
    Content       string           `json:"content" db:"content"`
    Language      string           `json:"language" db:"language"`
}

// Notification History Table
type NotificationHistory struct {
    ID                uuid.UUID `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID            uuid.UUID `json:"user_id" db:"user_id"`
    NotificationType  NotificationType `json:"notification_type" db:"notification_type"`
    Status            string    `json:"status" db:"status"`
    SentAt            *time.Time `json:"sent_at,omitempty" db:"sent_at"`
}

// Rate Limit Table
type RateLimit struct {
    ID               uuid.UUID `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID           uuid.UUID `json:"user_id" db:"user_id"`
    IPAddress        string    `json:"ip_address" db:"ip_address"`
    Requests         int       `json:"requests" db:"requests"`
    Limit            int       `json:"limit" db:"limit"`
    ResetAt          time.Time `json:"reset_at" db:"reset_at"`
}
