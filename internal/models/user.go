package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID          int       `json:"id" gorm:"primarykey"`
	Username    string    `json:"username" gorm:"column:username;type:varchar(20)" validate:"required"`
	Email       string    `json:"email" gorm:"column:email;type:varchar(100)" validate:"required"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number;type:varchar(15)"`
	FullName    string    `json:"full_name" gorm:"column:full_name;type:varchar(100)" validate:"required"`
	Address     string    `json:"address" gorm:"column:address;type:text"`
	DateOfBirth string    `json:"dob" gorm:"column:dob;type:date"`
	Password    string    `json:"password" gorm:"column:password;type:varchar(255)" validate:"required"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (*User) TableName() string {
	return "users"
}

func (m *User) Validate() error {
	v := validator.New()
	return v.Struct(m)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (m *LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(m)
}

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
