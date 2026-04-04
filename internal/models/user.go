package models

import "time"

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username" gorm:"column:username;type:varchar(20)"`
	Email       string    `json:"email" gorm:"column:email;type:varchar(100)"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone_number;type:varchar(15)"`
	FullName    string    `json:"full_name" gorm:"column:full_name;type:varchar(100)"`
	Address     string    `json:"address" gorm:"column:address;type:text"`
	DateOfBirth string    `json:"dob" gorm:"column:dob;type:date"`
	Password    string    `json:"password" gorm:"column:password;type:varchar(255)"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func (*User) TableName() string {
	return "users"
}

type UserSession struct {
	ID                  uint      `gorm:"primarykey"`
	UserID              uint      `json:"user_id" gorm:"type:int" validate:"required"`
	Token               string    `json:"token" gorm:"type:varchar(255)" validate:"required"`
	RefreshToken        string    `json:"refresh_token" gorm:"type:varchar(255)" validate:"required"`
	TokenExpired        time.Time `json:"-" validate:"required"`
	RefreshTokenExpired time.Time `json:"-" validate:"required"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (*UserSession) TableName() string {
	return "user_sessions"
}
