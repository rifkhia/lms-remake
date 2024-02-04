package utils

import (
	"github.com/rifkhia/lms-remake/internal/pkg"
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(userPassword string, inputPassword string) pkg.CustomError {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(inputPassword))
	if err != nil {
		return pkg.CustomError{
			Cause:   err,
			Service: inputPassword,
			Code:    BAD_REQUEST,
		}
	}
	return pkg.CustomError{}
}
