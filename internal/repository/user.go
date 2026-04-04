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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) int {
	var user models.User
	r.DB.Where("email = ?", email).First(&user)

	return user.ID
}

func (r *UserRepository) Register(ctx context.Context, user models.User) error {
	err := gorm.G[models.User](r.DB).Create(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}
