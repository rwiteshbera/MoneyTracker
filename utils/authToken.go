package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type SignedUserDetails struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	jwt.RegisteredClaims
}

// Generate JWT Token
func GenerateToken(firstname string, lastname string, phonenumber string) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		return "", err
	}

	SECRET_KEY := os.Getenv("SECRET_KEY")

	claims := SignedUserDetails{
		FirstName:   firstname,
		LastName:    lastname,
		PhoneNumber: phonenumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 168)), // 7 Days
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Validate JWT Token

func ValidateToken(authToken string) (*SignedUserDetails, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	SECRET_KEY := os.Getenv("SECRET_KEY")

	token, err := jwt.ParseWithClaims(
		authToken,
		&SignedUserDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedUserDetails)
	if !ok {
		return nil, errors.New("token is invalid")
	}
	return claims, nil
}
