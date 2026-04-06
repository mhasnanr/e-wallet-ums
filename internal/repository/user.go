package repository

import (
	"context"
	"time"

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
	err := r.DB.Where("email = ?", email).First(&user)
	if err.Error != nil {
		return user, err.Error
	}
	return user, nil
}

func (r *UserRepository) Register(ctx context.Context, user models.User) error {
	err := gorm.G[models.User](r.DB).Create(ctx, &user)
	if err != nil {
		return err
	}

	return nil
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
	err := r.DB.Where("refresh_token = ?", refreshToken).First(&userSession)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *UserRepository) UpdateTokenByRefreshToken(ctx context.Context, token string, refreshToken string) error {
	var userSession models.UserSession
	err := r.DB.Where("refresh_token = ?", refreshToken).First(&userSession).Error
	if err != nil {
		return err
	}

	err = r.DB.Model(&userSession).Update("token", token).Updates(map[string]interface{}{
		"token":         token,
		"token_expired": time.Now().Add(helpers.MapTypeToken["token"]),
	}).Error
	if err != nil {
		return err
	}

	return nil
}
