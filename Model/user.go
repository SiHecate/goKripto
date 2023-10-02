package model

import (
	"gorm.io/gorm"
)

// User model user information.
type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password []byte `json:"-"`
	Wallet   Wallet `json:"wallet" gorm:"foreignKey:UserID"`
	Deneme   bool   `json:"boolean"`
}
