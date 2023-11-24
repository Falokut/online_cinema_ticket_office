package model

import (
	"time"
)

type Account struct {
	UUID             string    `db:"id"`
	Email            string    `db:"email"`
	Password         string    `db:"password_hash"`
	RegistrationDate time.Time `db:"registration_date"`
}

type UserProfile struct {
	Username         string    `db:"username"`
	Email            string    `db:"email"`
	Password         string    `db:"password_hash"`
	RegistrationDate time.Time `db:"registration_date"`
}
