package services

import (
	"context"
	"errors"
	"time"

	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
)

type JWTManager interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashed string, plain string) error
	GenerateToken(user models.User, tokenType string) (string, error)
	ValidateToken(ctx context.Context, token string) (*helpers.ClaimToken, error)
}

type UserRepository interface {
	Register(context.Context, models.User) error
	GetUserByEmail(context.Context, string) (models.User, error)
	CreateUserSession(context.Context, models.UserSession) error
	GetUserSessionByRefreshToken(context.Context, string) error
	UpdateTokenByRefreshToken(context.Context, string, string) error
}

type UserService struct {
	repo       UserRepository
	jwtManager JWTManager
}

func NewUserService(repo UserRepository, jwtManager JWTManager) *UserService {
	return &UserService{repo: repo, jwtManager: jwtManager}
}

func (s *UserService) Register(ctx context.Context, user models.User) error {
	returnedUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if returnedUser.ID != 0 {
		return errors.New(constants.ErrDuplicateEmail)
	}

	if err != nil {
		return err
	}

	hashedPassword, err := s.jwtManager.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return s.repo.Register(ctx, user)
}

func (s *UserService) Login(ctx context.Context, req models.LoginRequest) (models.LoginResponse, error) {
	var (
		response models.LoginResponse
		now      = time.Now()
	)

	returnedUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if returnedUser.ID == 0 {
		return response, errors.New(constants.ErrUserNotFound)
	}

	if err != nil {
		return response, err
	}

	err = s.jwtManager.VerifyPassword(returnedUser.Password, req.Password)
	if err != nil {
		return response, err
	}

	token, err := s.jwtManager.GenerateToken(returnedUser, "token")
	if err != nil {
		return response, errors.New("failed to generate token")
	}

	refreshToken, err := s.jwtManager.GenerateToken(returnedUser, "refreshToken")
	if err != nil {
		return response, errors.New("failed to generate token")
	}

	response.Token = token

	userSession := models.UserSession{
		UserID:              uint(returnedUser.ID),
		Token:               token,
		RefreshToken:        refreshToken,
		TokenExpired:        now.Add(helpers.MapTypeToken["token"]),
		RefreshTokenExpired: now.Add(helpers.MapTypeToken["refresh_token"]),
	}

	err = s.repo.CreateUserSession(ctx, userSession)
	if err != nil {
		return response, errors.New("failed to create user session")
	}

	response.UserID = returnedUser.ID
	response.Token = token
	response.RefreshToken = refreshToken

	return response, nil
}

func (s *UserService) UpdateTokenByRefreshToken(ctx context.Context, refreshToken string, claims helpers.ClaimToken) (string, error) {
	var user models.User

	user.ID = claims.UserID
	user.Username = claims.Username
	user.Email = claims.Email
	user.FullName = claims.Fullname

	newToken, err := s.jwtManager.GenerateToken(user, "token")
	if err != nil {
		return newToken, errors.New("failed to generate token")
	}

	err = s.repo.UpdateTokenByRefreshToken(ctx, newToken, refreshToken)
	if err != nil {
		return newToken, errors.New("failed to update token")
	}

	return newToken, nil
}
