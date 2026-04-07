package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"github.com/mhasnanr/ewallet-ums/internal/models"
)

type JWTManager struct{}

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

func (j *JWTManager) GenerateToken(user models.User, tokenType string) (string, error) {
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

func (j *JWTManager) ValidateToken(ctx context.Context, token string) (*ClaimToken, error) {
	var (
		claimToken *ClaimToken
		ok         bool
	)

	jwtToken, err := jwt.ParseWithClaims(token, &ClaimToken{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("failed to validate method jwt: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt: %v", err)
	}

	if claimToken, ok = jwtToken.Claims.(*ClaimToken); !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	return claimToken, nil
}
