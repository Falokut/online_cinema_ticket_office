package user_service

import "time"

type User struct {
	UUID             string    `json:"uuid" bson:"_id" db:"id"`
	Username         string    `json:"username" bson:"username" db:"username"`
	Email            string    `json:"email" bson:"email" db:"email"`
	Password         string    `json:"-" bson:"password,omitempty" db:"password_hash"`
	Verified         bool      `json:"verified" bson:"verified" db:"verified"`
	ProfilePictureID string    `json:"-" db:"profile_picture_id"`
	RegistrationDate time.Time `json:"registration_date" bson:"registration_date" db:"registration_date"`
}

type UpdateUserDTO struct {
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
	Verified    bool   `json:"verified" bson:"verified, omitempty"`
}

type UserProfile struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	RegistrationDate  time.Time `json:"registration_date"`
}

type userRequestCtx struct {
	AccessToken string `json:"token" binding:"required"`
}
