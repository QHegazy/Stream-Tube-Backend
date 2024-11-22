package models

type UserStatus string
type OAuthProvider string
type AuthMethod string
type MFAType string
type NotificationType string
type RiskLevel string
type Gender string

const (
	UserStatusActive       UserStatus = "active"
	UserStatusInactive     UserStatus = "inactive"
	UserStatusSuspended    UserStatus = "suspended"
	UserStatusBanned       UserStatus = "banned"
	UserStatusPendingVerification UserStatus = "pending_verification"

	OAuthProviderGoogle OAuthProvider = "google"
	OAuthProviderFacebook OAuthProvider = "facebook"
	OAuthProviderGitHub OAuthProvider = "github"
	OAuthProviderMicrosoft OAuthProvider = "microsoft"

	AuthMethodLocal AuthMethod = "local"
	AuthMethodOAuth AuthMethod = "oauth"

	MFATypeSMS MFAType = "sms"
	MFATypeEmail MFAType = "email"
	MFATypeTOTP MFAType = "totp"

	NotificationTypeSecurity NotificationType = "security"
	NotificationTypeAccount NotificationType = "account"
	NotificationTypeMarketing NotificationType = "marketing"
	NotificationTypeSystem NotificationType = "system"

	RiskLevelLow RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"

	GenderMale Gender = "male"
	GenderFemale Gender = "female"
	GenderOther Gender = "other"
)