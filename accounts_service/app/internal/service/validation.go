package service

import (
	"errors"
	"regexp"

	accounts_service "github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/accounts_service/v1/protos"
)

type validateError struct {
	DeveloperMessage string
	UserMessage      string
}

func newValidateError(DeveloperMessage, UserMessage string) *validateError {
	return &validateError{DeveloperMessage: DeveloperMessage, UserMessage: UserMessage}
}
func (e *validateError) Error() string {
	err := errors.Join(errors.New(e.DeveloperMessage), errors.New(e.UserMessage))
	return err.Error()
}

var (
	ErrInputBodyNotValid  = errors.New("input body not valid")
	ErrPasswordsDontMatch = errors.New("passwords don't match")
)

func validateSignupInput(input *accounts_service.CreateAccountRequest) *validateError {

	if input == nil {
		return newValidateError(ErrInputBodyNotValid.Error(), "")
	}

	if input.Password != input.RepeatPassword {
		return newValidateError("", ErrPasswordsDontMatch.Error())
	}

	if err := validatePassword(input.Password); err != nil {
		return err
	}
	if err := validateEmail(input.Email); err != nil {
		return err
	}

	return nil
}

func validatePassword(Password string) *validateError {
	passwordLengh := len(Password)
	if passwordLengh < 6 || passwordLengh > 32 {
		return newValidateError("", "the password must be less than 32 symbols and more than 6")
	}

	return nil
}

func validateEmail(email string) *validateError {
	if len(email) > 100 || len(email) < 4 {
		return newValidateError("", "email must be less than 100 symbols and more than 4")
	}
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(email) {
		return newValidateError("", "email is not valid")
	}
	return nil
}
