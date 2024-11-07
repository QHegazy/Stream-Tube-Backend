package models


type AuthMethod string

const (
    AuthMethodLocal AuthMethod = "local"
    AuthMethodOAuth AuthMethod = "oauth"
    AuthMethodMFA   AuthMethod = "mfa"
)


type UserStatus string
const (
    Active               UserStatus = "active"
    Inactive             UserStatus = "inactive"
    Suspended            UserStatus = "suspended"
    Banned               UserStatus = "banned"
    PendingVerification  UserStatus = "pending_verification"
)

type OAuthProvider string
const (
    ProviderGoogle    OAuthProvider = "google"
    ProviderFacebook  OAuthProvider = "facebook"
    ProviderGithub    OAuthProvider = "github"
    ProviderMicrosoft OAuthProvider = "microsoft"
)

type MFAType string
const (
    Authenticator MFAType = "authenticator"
    SMS           MFAType = "sms"
    Email         MFAType = "email"
    SecurityKey   MFAType = "security_key"
)

type NotificationType string
const (
    Security   NotificationType = "security"
    Account    NotificationType = "account"
    Marketing  NotificationType = "marketing"
    System     NotificationType = "system"
)

type RiskLevel string
const (
    Low      RiskLevel = "low"
    Medium   RiskLevel = "medium"
    High     RiskLevel = "high"
    Critical RiskLevel = "critical"
)