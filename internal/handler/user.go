package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
	"github.com/mhasnanr/ewallet-ums/internal/models/dto"
)

type UserService interface {
	Register(context.Context, *models.User) (*models.User, error)
	Login(context.Context, models.LoginRequest) (dto.LoginResponse, error)
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
		r.writeErrorResponse(c, constants.ErrorBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		r.writeErrorResponse(c, err, nil)
		return
	}

	user, err := r.service.Register(c.Request.Context(), &req)
	if err != nil {
		r.writeErrorResponse(c, err, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusCreated, constants.UserCreated, user)
}

func (r *UserHandler) login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		r.writeErrorResponse(c, constants.ErrorBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		r.writeErrorResponse(c, err, nil)
		return
	}

	res, err := r.service.Login(c.Request.Context(), req)
	if err != nil {
		r.writeErrorResponse(c, err, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusOK, constants.LoginSucceed, res)
}

func (r *UserHandler) logout(c *gin.Context) {
	token, ok := c.Get("accessToken")
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToGetToken, nil)
		return
	}

	accessToken, ok := token.(string)
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToParseToken, nil)
		return
	}

	err := r.service.Logout(c.Request.Context(), accessToken)
	if err != nil {
		r.writeErrorResponse(c, err, nil)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusNoContent, constants.LoginSucceed, nil)
}

func (r *UserHandler) refreshToken(c *gin.Context) {
	var response dto.RefreshTokenResponse

	token, ok := c.Get("refreshToken")
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToGetToken, nil)
		return
	}

	refreshToken, ok := token.(string)
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToParseToken, nil)
		return
	}

	val, ok := c.Get("claim")
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToExtractClaims, nil)
		return
	}

	claim, ok := val.(*helpers.ClaimToken)
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToParseClaims, nil)
		return
	}

	newToken, err := r.service.UpdateTokenByRefreshToken(c.Request.Context(), refreshToken, claim)
	if err != nil {
		r.writeErrorResponse(c, constants.ErrorFailedToUpdateToken, nil)
		return
	}

	response.Token = newToken

	helpers.SendResponseHTTP(c, http.StatusOK, constants.NewToken, response)
}

func (r *UserHandler) validateToken(c *gin.Context) {
	var response dto.ValidateTokenResponse

	token, ok := c.Get("accessToken")
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToGetToken, nil)
		return
	}

	_, ok = token.(string)
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToParseToken, nil)
		return
	}

	val, ok := c.Get("claim")
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToExtractClaims, nil)
		return
	}

	claim, ok := val.(*helpers.ClaimToken)
	if !ok {
		r.writeErrorResponse(c, constants.ErrorFailedToParseClaims, nil)
		return
	}

	response.UserID = claim.UserID

	helpers.SendResponseHTTP(c, http.StatusOK, constants.ValidToken, response)
}

func (r *UserHandler) writeErrorResponse(c *gin.Context, err error, data any) {
	var appErr *constants.AppError
	var valErrs validator.ValidationErrors

	if errors.As(err, &appErr) {
		helpers.SendResponseHTTP(c, appErr.StatusCode, appErr.Message, data)
		return
	}

	if errors.As(err, &valErrs) {
		errStr := helpers.ConstructErrString(valErrs)
		helpers.SendResponseHTTP(c, http.StatusBadRequest, errStr, data)
		return
	}

	helpers.SendResponseHTTP(c, http.StatusInternalServerError, err.Error(), nil)
}
