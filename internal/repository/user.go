package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mhasnanr/ewallet-ums/constants"
	"github.com/mhasnanr/ewallet-ums/helpers"
	"github.com/mhasnanr/ewallet-ums/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	res := r.DB.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return user, res.Error
	}
	return user, nil
}

func (r *UserRepository) Register(ctx context.Context, user *models.User) (*models.User, error) {
	err := gorm.G[models.User](r.DB).Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateUserSession(ctx context.Context, userSession models.UserSession) error {
	err := gorm.G[models.UserSession](r.DB).Create(ctx, &userSession)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserSessionByRefreshToken(ctx context.Context, refreshToken string) error {
	var userSession models.UserSession
	err := r.DB.Where("refresh_token = ?", refreshToken).First(&userSession).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrorSessionNotFound
		}
		return err
	}
	return nil
}

func (r *UserRepository) UpdateTokenByRefreshToken(ctx context.Context, token string, refreshToken string) error {
	err := r.DB.Model(&models.UserSession{}).
		Where("refresh_token = ?", refreshToken).
		Updates(map[string]any{
			"token":         token,
			"token_expired": time.Now().Add(helpers.MapTypeToken["token"]),
		}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrorSessionNotFound
		}
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userID int) error {
	err := r.DB.Where("id = ?", userID).Delete(&models.User{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrorSessionNotFound
		}
		return err
	}

	return nil
}
