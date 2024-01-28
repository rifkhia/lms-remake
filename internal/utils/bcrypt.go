package utils

import "golang.org/x/crypto/bcrypt"

func GeneratePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(studentPassword string, inputPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(studentPassword), []byte(inputPassword))
	if err != nil {
		return err
	}
	return nil
}
