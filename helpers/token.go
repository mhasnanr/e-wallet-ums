package helpers

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type JWTManager interface {
	HashPassword(password string) (string, error)
}

type JWTApp struct{}

func (j *JWTApp) HashPassword(password string) (string, error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(hashedByte), nil
}
