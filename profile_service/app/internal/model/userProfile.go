package model

import (
	"time"
)

type UserProfile struct {
	Username         string    `db:"username"`
	Email            string    `db:"email"`
	ProfilePictureID string    `db:"profile_picture_id"`
	RegistrationDate time.Time `db:"registration_date"`
}
