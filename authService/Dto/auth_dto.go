package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username       string    `json:"username" validate:"required,min=3,max=50,alphanumunicode"`
	Email          string    `json:"email" validate:"required,email"`
	Password       string    `json:"password" validate:"required,min=8"`
	FullName       string    `json:"full_name"`
	BirthDate      time.Time `json:"birth_date"`
	Gender         string    `json:"gender,omitempty"`
	Location       string    `json:"location,omitempty"`
	Timezone       string    `json:"timezone,omitempty"`
	Language       string    `json:"language"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID             uuid.UUID     `json:"id"`
	Username       string        `json:"username"`
	Email          string        `json:"email"`
	Status         string        `json:"status"`
	EmailVerified  *time.Time    `json:"email_verified,omitempty"`
	LastActiveAt   *time.Time    `json:"last_active_at,omitempty"`
	Profile        *ProfileResponse `json:"profile,omitempty"`
	OAuthProviders []OAuthProvider `json:"oauth_providers,omitempty"`
}

type ProfileResponse struct {
	ID                uuid.UUID  `json:"id"`
	FullName          string     `json:"full_name"`
	ProfilePictureURL string     `json:"profile_picture_url,omitempty"`
	CoverPhotoURL     string     `json:"cover_photo_url,omitempty"`
	BirthDate         time.Time  `json:"birth_date"`
	Gender            string     `json:"gender,omitempty"`
	Location          string     `json:"location,omitempty"`
	Timezone          string     `json:"timezone,omitempty"`
	Language          string     `json:"language"`
}

type OAuthProvider string

const (
	GoogleOAuthProvider OAuthProvider = "google"
	MicrosoftOAuthProvider OAuthProvider = "microsoftonline"
)
