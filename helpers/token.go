package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"github.com/mhasnanr/ewallet-ums/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type JWTApp struct{}

type ClaimToken struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

var MapTypeToken = map[string]time.Duration{
	"token":        time.Hour * 3,
	"refreshToken": time.Hour * 72,
}

var jwtSecret = []byte(bootstrap.GetEnv("APP_SECRET", ""))

func (j *JWTApp) HashPassword(password string) (string, error) {
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password")
	}

	return string(hashedByte), nil
}

func (j *JWTApp) VerifyPassword(hashed string, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	if err != nil {
		return fmt.Errorf("password is invalid")
	}

	return nil
}

func (j *JWTApp) GenerateToken(user models.User, tokenType string) (string, error) {
	now := time.Now()

	claimToken := ClaimToken{
		UserID:   user.ID,
		Username: user.Username,
		Fullname: user.FullName,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    bootstrap.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(MapTypeToken[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return tokenString, fmt.Errorf("failed to generate token: %v", err)
	}

	return tokenString, nil
}
