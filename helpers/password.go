package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct{}

func (p *PasswordHasher) HashPassword(password string) (string, error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")
	}

	return string(hashedByte), nil
}

func (p *PasswordHasher) VerifyPassword(hashed string, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return fmt.Errorf("password is invalid")
	}

	return nil
}
