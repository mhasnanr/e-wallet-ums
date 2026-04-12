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

type SessionRepository struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) CreateUserSession(ctx context.Context, userSession models.UserSession) error {
	err := gorm.G[models.UserSession](r.DB).Create(ctx, &userSession)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) GetUserSessionByRefreshToken(ctx context.Context, refreshToken string) error {
	var userSession models.UserSession
	err := r.DB.Where("refresh_token = ?", refreshToken).Last(&userSession)
	if err.Error != nil {
		return err.Error
	}
	return nil
}

func (r *SessionRepository) UpdateTokenByRefreshToken(ctx context.Context, token string, refreshToken string) error {
	var userSession models.UserSession
	err := r.DB.Where("refresh_token = ?", refreshToken).Last(&userSession).Error
	if err != nil {
		return err
	}

	err = r.DB.Model(&userSession).Update("token", token).Updates(map[string]any{
		"token":         token,
		"token_expired": time.Now().Add(helpers.MapTypeToken["token"]),
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) DeleteUserSessionByToken(ctx context.Context, accessToken string) error {
	err := r.DB.Where("token = ?", accessToken).Delete(&models.UserSession{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrorUserNotFound
		}

		return err
	}

	return nil
}
