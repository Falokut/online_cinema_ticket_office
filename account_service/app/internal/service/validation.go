package service

import (
	"errors"
	"regexp"

	account_service "github.com/Falokut/online_cinema_ticket_office/account_service/pkg/account_service/protos"
)

func validateSignupInput(input *account_service.CreateAccountRequest) error {

	if input == nil {
		return errors.New("Input body not valid.")
	}

	if input.Password != input.RepeatPassword {
		return errors.New("Passwords don't match.")
	}

	if err := validatePassword(input.Password); err != nil {
		return err
	}
	if err := validateEmail(input.Email); err != nil {
		return err
	}

	return nil
}

func stringContainsChar(str []byte, toFind byte) bool {
	for _, char := range str {
		if char == toFind {
			return true
		}
	}
	return false
}

func validatePassword(Password string) error {
	passwordLengh := len(Password)
	if passwordLengh < 6 || passwordLengh > 32 {
		return errors.New("The password must be less than 32 symbols and more than 6.")
	}

	return nil
}

func validateEmail(email string) error {
	if len(email) > 100 || len(email) < 4 {
		return errors.New("Email must be less than 100 symbols and more than 4.")
	}
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(email) {
		return errors.New("Email is not valid")
	}
	return nil
}
