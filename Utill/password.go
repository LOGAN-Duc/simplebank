package Utill

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassWord return the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hash password: %w", err)
	}
	return string(HashPassword), nil
}

//CheckPassword

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
