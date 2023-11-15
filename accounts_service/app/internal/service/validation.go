package service

import (
	"errors"
	"regexp"

	accounts_service "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/accounts_service/v1/protos"
)

func validateSignupInput(input *accounts_service.CreateAccountRequest) error {

	if input == nil {
		return errors.New("input body not valid")
	}

	if input.Password != input.RepeatPassword {
		return errors.New("passwords don't match")
	}

	if err := validatePassword(input.Password); err != nil {
		return err
	}
	if err := validateEmail(input.Email); err != nil {
		return err
	}

	return nil
}

func validatePassword(Password string) error {
	passwordLengh := len(Password)
	if passwordLengh < 6 || passwordLengh > 32 {
		return errors.New("the password must be less than 32 symbols and more than 6")
	}

	return nil
}

func validateEmail(email string) error {
	if len(email) > 100 || len(email) < 4 {
		return errors.New("email must be less than 100 symbols and more than 4")
	}
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(email) {
		return errors.New("email is not valid")
	}
	return nil
}
