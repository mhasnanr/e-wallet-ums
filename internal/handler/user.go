package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
)

type UserService interface {
	Register(context.Context, models.User) error
	Login(context.Context, models.LoginRequest) (models.LoginResponse, error)
	UpdateTokenByRefreshToken(context.Context, string, helpers.ClaimToken) (string, error)
}

type AuthMiddleware interface {
	MiddlewareRefreshToken(c *gin.Context)
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
	userV1.GET("/refresh-token", r.authMiddleware.MiddlewareRefreshToken, r.refreshToken)

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

func (r *UserHandler) refreshToken(c *gin.Context) {
	token, ok := c.Get("refreshToken")
	if !ok {
		fmt.Printf("failed to get refresh token")
	}

	refreshToken, ok := token.(string)
	if !ok {
		fmt.Printf("failed to parse refresh token")
	}

	val, ok := c.Get("claim")
	if !ok {
		fmt.Printf("failed to get token claim")
	}

	claim, ok := val.(helpers.ClaimToken)
	if !ok {
		fmt.Printf("failed to parse token claim")
	}

	newToken, err := r.service.UpdateTokenByRefreshToken(c.Request.Context(), refreshToken, claim)
	if err != nil {
		helpers.SendResponseHTTP(c, http.StatusOK, constants.MsgNewToken, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.MsgNewToken, map[string]any{
		"token": newToken,
	})
}
