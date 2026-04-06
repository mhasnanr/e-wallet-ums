package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/bootstrap"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
)

type UserRepository interface {
	Register(context.Context, models.User) error
	GetUserByEmail(context.Context, string) (models.User, error)
	CreateUserSession(context.Context, models.UserSession) error
	GetUserSessionByRefreshToken(context.Context, string) error
}

type JWTManager interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashed string, plain string) error
	GenerateToken(user models.User, tokenType string) (string, error)
	ValidateToken(ctx context.Context, token string) (*helpers.ClaimToken, error)
}

type AuthMiddleware struct {
	repo       UserRepository
	jwtManager JWTManager
}

func NewAuthMiddleware(repo UserRepository, jwtManager JWTManager) *AuthMiddleware {
	return &AuthMiddleware{repo, jwtManager}
}

func (a *AuthMiddleware) MiddlewareRefreshToken(c *gin.Context) {
	var log = bootstrap.Log

	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		log.Infow("authorization empty")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	refreshToken := strings.Split(auth, "Bearer ")[1]
	if refreshToken == "" {
		log.Infow("invalid token")
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	err := a.repo.GetUserSessionByRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		log.Infow("failed to get user session on DB: ", err)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	claim, err := a.jwtManager.ValidateToken(c.Request.Context(), refreshToken)
	if err != nil {
		log.Infow(err.Error())
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Infow("jwt token is expired: ", claim.ExpiresAt)
		helpers.SendResponseHTTP(c, http.StatusUnauthorized, "unauthorized", nil)
		c.Abort()
		return
	}

	c.Set("refreshToken", refreshToken)
	c.Set("claim", claim)
	c.Next()
}
