package services

import (
	"context"
	"errors"

	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
)

type UserRepository interface {
	Register(context.Context, models.User) error
	GetUserByEmail(context.Context, string) int
}

type UserService struct {
	repo       UserRepository
	jwtManager helpers.JWTManager
}

func NewUserService(repo UserRepository, jwtManager helpers.JWTManager) *UserService {
	return &UserService{repo: repo, jwtManager: jwtManager}
}

func (s *UserService) Register(ctx context.Context, user models.User) error {
	userID := s.repo.GetUserByEmail(ctx, user.Email)
	if userID != 0 {
		return errors.New(constants.ErrDuplicateEmail)
	}

	hashedPassword, err := s.jwtManager.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return s.repo.Register(ctx, user)
}
