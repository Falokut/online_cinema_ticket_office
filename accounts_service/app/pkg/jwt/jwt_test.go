package jwt_test

import (
	"testing"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/accounts_service/pkg/jwt"
)

func TestJWT(t *testing.T) {
	testCases := []struct {
		Value    string
		Secret   string
		TokenTTL time.Duration
	}{
		{
			Value:    "sdakjdjakldad",
			Secret:   "asdkwwq",
			TokenTTL: time.Hour,
		},
		{
			Value:    "sadasdasd",
			Secret:   "asdadasda",
			TokenTTL: time.Hour,
		},
		{
			Value:    "asdasdas",
			Secret:   "sada",
			TokenTTL: time.Hour,
		},
	}

	Tokens := make(map[string]struct {
		Secret              string
		ExpectedParsedValue string
	})

	for _, testCase := range testCases {
		token, err := jwt.GenerateToken(testCase.Value, testCase.Secret, testCase.TokenTTL)
		if err != nil {
			t.Errorf("Something wrong, getting error:%s", err.Error())
		}
		Tokens[token] = struct {
			Secret              string
			ExpectedParsedValue string
		}{testCase.Secret, testCase.Value}
	}

	for token, values := range Tokens {
		parsedValue, err := jwt.ParseToken(token, values.Secret)
		if err != nil {
			t.Errorf("Something wrong, getting error:%s", err.Error())
		}

		if values.ExpectedParsedValue != parsedValue {
			t.Errorf("Result was incorrect, got %s , want %s", parsedValue, values.ExpectedParsedValue)
		}
	}
}
