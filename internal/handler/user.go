package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
)

type RegisterService interface {
	Register(context.Context, models.User) error
	Login(context.Context, models.LoginRequest) (models.LoginResponse, error)
}

type UserHandler struct {
	service RegisterService
}

func NewUserHandler(svc RegisterService) *UserHandler {
	return &UserHandler{service: svc}
}

func (r *UserHandler) RegisterRoute(c *gin.Engine) {
	userV1 := c.Group("/users/v1")
	userV1.POST("/register", r.registerUser)
	userV1.POST("/login", r.login)

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
