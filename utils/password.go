package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// Encrypt the password
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("unable to hash the password")
	}
	return string(hashed), nil
}

// Verify the user password with saved user password
func VerifyPassword(userPassword string, hashedPassword string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword)); err != nil {
		return false, errors.New("wrong credentials")
	}
	return true, nil
}
