package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type UserSession struct {
	ID                  uint      `gorm:"primarykey"`
	UserID              uint      `json:"user_id" gorm:"type:int" validate:"required"`
	Token               string    `json:"token" gorm:"type:text" validate:"required"`
	RefreshToken        string    `json:"refresh_token" gorm:"type:text" validate:"required"`
	TokenExpired        time.Time `json:"-" validate:"required"`
	RefreshTokenExpired time.Time `json:"-" validate:"required"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (m *UserSession) Validate() error {
	v := validator.New()
	return v.Struct(m)
}

func (*UserSession) TableName() string {
	return "user_sessions"
}
