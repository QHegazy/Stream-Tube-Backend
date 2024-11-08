package models

import (
	"time"

	"github.com/google/uuid"
)


type Profile struct {
    ID                uuid.UUID  `json:"id" db:"id" default:"uuid_generate_v4()"`
    UserID            uuid.UUID  `json:"user_id" db:"user_id"`
    FullName          string     `json:"full_name" db:"full_name"`
    ProfilePictureURL string     `json:"profile_picture_url,omitempty" db:"profile_picture_url"`
    CoverPhotoURL     string     `json:"cover_photo_url,omitempty" db:"cover_photo_url"`
    BirthDate         time.Time  `json:"birth_date" db:"birth_date"`
    Gender            string     `json:"gender,omitempty" db:"gender"`
    Location          string     `json:"location,omitempty" db:"location"`
    Language          string     `json:"language" db:"language" default:"en"`
}
