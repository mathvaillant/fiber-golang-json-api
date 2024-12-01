package main

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateToken(user *User) (string, error) {
	secret := []byte("super-secret-key")
	method := jwt.SigningMethodHS256
	claims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 168).Unix(), // 7 days
	}

	token, err := jwt.NewWithClaims(method, claims).SignedString(secret)
	if err != nil {
		return "", err
	}

	return token, nil
}
