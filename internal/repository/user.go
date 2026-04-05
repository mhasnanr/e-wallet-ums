package repository

import (
	"context"

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
	tx := r.DB.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		return user, tx.Error
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
