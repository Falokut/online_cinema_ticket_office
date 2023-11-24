package model

import (
	"database/sql"
	"time"
)

type UserProfile struct {
	AccountID        string         `db:"account_id"`
	Username         string         `db:"username"`
	Email            string         `db:"email"`
	ProfilePictureID sql.NullString `db:"profile_picture_id"`
	RegistrationDate time.Time      `db:"registration_date"`
}
