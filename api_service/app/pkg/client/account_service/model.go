package account_service

type SignupUserDTO struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RepeatPassword string `json:"repeat_password" binding:"required"`
}

type SigninUserDTO struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password,omitempty"`
}
