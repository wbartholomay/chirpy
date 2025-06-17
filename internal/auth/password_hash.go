package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(([]byte)(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error encrypting password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword(([]byte)(hash), ([]byte)(password))
}