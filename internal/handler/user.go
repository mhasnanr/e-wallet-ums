package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
)

type UserService interface {
	Register(context.Context, models.User) error
	Login(context.Context, models.LoginRequest) (models.LoginResponse, error)
	UpdateTokenByRefreshToken(context.Context, string, *helpers.ClaimToken) (string, error)
	Logout(context.Context, string) error
}

type AuthMiddleware interface {
	MiddlewareRefreshToken(c *gin.Context)
	MiddlewareAccessToken(c *gin.Context)
}

type UserHandler struct {
	service        UserService
	authMiddleware AuthMiddleware
}

func NewUserHandler(svc UserService, authMiddleware AuthMiddleware) *UserHandler {
	return &UserHandler{service: svc, authMiddleware: authMiddleware}
}

func (r *UserHandler) RegisterRoute(c *gin.Engine) {
	userV1 := c.Group("/users/v1")
	userV1.POST("/register", r.registerUser)
	userV1.POST("/login", r.login)
	userV1.DELETE("/logout", r.authMiddleware.MiddlewareAccessToken, r.logout)
	userV1.GET("/token/refresh", r.authMiddleware.MiddlewareRefreshToken, r.refreshToken)
	userV1.GET("/token/validate", r.authMiddleware.MiddlewareAccessToken, r.validateToken)

}

func (r *UserHandler) registerUser(c *gin.Context) {
	var req models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFieldBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFieldBadRequest, nil)
		return
	}

	err := r.service.Register(c.Request.Context(), req)
	if err != nil {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusCreated, constants.MsgUserCreated, nil)
}

func (r *UserHandler) login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFieldBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		helpers.SendResponseHTTP(c, http.StatusBadRequest, constants.ErrFieldBadRequest, nil)
		return
	}

	res, err := r.service.Login(c.Request.Context(), req)
	if err != nil {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusCreated, constants.MsgLoginSucceed, res)
}

func (r *UserHandler) logout(c *gin.Context) {
	token, ok := c.Get("accessToken")
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "failed to get token", nil)
		return
	}

	accessToken, ok := token.(string)
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "failed to parse token", nil)
		return
	}

	err := r.service.Logout(c.Request.Context(), accessToken)
	if err != nil {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusNoContent, constants.MsgLogoutSucceed, nil)
}

func (r *UserHandler) refreshToken(c *gin.Context) {
	token, ok := c.Get("refreshToken")
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "failed to get refresh token", nil)
		return
	}

	refreshToken, ok := token.(string)
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "failed to parse token", nil)
		return
	}

	val, ok := c.Get("claim")
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "user data is not found", nil)
		return
	}

	claim, ok := val.(*helpers.ClaimToken)
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "failed to parse user data", nil)
		return
	}

	newToken, err := r.service.UpdateTokenByRefreshToken(c.Request.Context(), refreshToken, claim)
	if err != nil {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "failed to update access token", nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.MsgNewToken, map[string]any{
		"token": newToken,
	})
}

func (r *UserHandler) validateToken(c *gin.Context) {
	token, ok := c.Get("accessToken")
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "failed to get access token", nil)
		return
	}

	_, ok = token.(string)
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "access token must be string", nil)
		return
	}

	val, ok := c.Get("claim")
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "user data is not found", nil)
		return
	}

	claim, ok := val.(*helpers.ClaimToken)
	if !ok {
		helpers.SendResponseHTTP(c, http.StatusInternalServerError, "failed to parse user data", nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, "token is valid", map[string]any{
		"user_id": claim.UserID,
	})
}
