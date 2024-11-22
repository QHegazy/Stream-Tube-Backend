package models

import (
	"time"

	"github.com/google/uuid"
)


type Profile struct {
	ID              uuid.UUID      `json:"id" db:"id"`
	UserID          uuid.UUID      `json:"user_id" db:"user_id"`
	ProfilePictureURL string        `json:"profile_picture_url" db:"profile_picture_url"`
	CoverPhotoURL   string         `json:"cover_photo_url" db:"cover_photo_url"`
	BirthDate       time.Time      `json:"birth_date" db:"birth_date"`
	Gender          Gender         `json:"gender" db:"gender"`
	Language        string         `json:"language" db:"language"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt       *time.Time     `json:"deleted_at" db:"deleted_at"`
}

type Address struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	AddressLine1  string         `json:"address_line1" db:"address_line1"`
	AddressLine2  string         `json:"address_line2" db:"address_line2"`
	City          string         `json:"city" db:"city"`
	Country       string         `json:"country" db:"country"`
	Timezone      string         `json:"timezone" db:"timezone"`
}

type SecondaryEmail struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	Email        string         `json:"email" db:"email"`
	VerifiedAt   time.Time      `json:"verified_at" db:"verified_at"`
}

type EmergencyContact struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	ContactName  string         `json:"contact_name" db:"contact_name"`
	ContactPhone string         `json:"contact_phone" db:"contact_phone"`
}

type Contact struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	UserID        uuid.UUID      `json:"user_id" db:"user_id"`
	PhoneNumber   string         `json:"phone_number" db:"phone_number"`
	AddressID     uuid.UUID      `json:"address_id" db:"address_id"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time     `json:"deleted_at" db:"deleted_at"`
}

type SecondaryContact struct {
	ID                  uuid.UUID      `json:"id" db:"id"`
	UserID              uuid.UUID      `json:"user_id" db:"user_id"`
	SecondaryEmailID    uuid.UUID      `json:"secondary_email_id" db:"secondary_email_id"`
	EmergencyContactID  uuid.UUID      `json:"emergency_contact_id" db:"emergency_contact_id"`
	CreatedAt           time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at" db:"updated_at"`
}