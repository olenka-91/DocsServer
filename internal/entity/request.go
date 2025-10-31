package entity

import (
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type SignUpRequest struct {
	Name     string `json:"name" binding:"required,alphanum,min=8,max=100"`
	Password string `json:"password" binding:"required,password_complexity,min=8"`
}

type SignInRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func ValidatePasswordComplexity(fl validator.FieldLevel) bool {
	pswd := fl.Field().String()

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	for _, ch := range pswd {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case !unicode.IsLetter(ch) && !unicode.IsDigit(ch):
			hasSpecial = true
		}
	}
	logrus.Info("hasUpper=", hasUpper, " hasLower=", hasLower, " hasDigit=", hasDigit, " hasSpecial=", hasSpecial)
	return hasUpper && hasLower && hasDigit && hasSpecial
}
