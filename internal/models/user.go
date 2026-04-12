package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID          int       `json:"id" gorm:"primarykey"`
	Username    string    `json:"username" gorm:"column:username;type:varchar(20);not null" validate:"required"`
	FullName    string    `json:"full_name" gorm:"column:full_name;type:varchar(100);not null" validate:"required"`
	Email       string    `json:"email" gorm:"column:email;type:varchar(100);not null" validate:"required"`
	Password    string    `json:"password,omitempty" gorm:"column:password;type:varchar(255);not null" validate:"required"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number;type:varchar(15)"`
	Address     string    `json:"address" gorm:"column:address;type:text"`
	DateOfBirth string    `json:"dob" gorm:"column:dob;type:date"`
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

func (m *LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(m)
}
