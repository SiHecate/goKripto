package model

import (
	"gorm.io/gorm"
)

// User model user information.
type User struct {
	gorm.Model
	Name        string `json:"name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Password    []byte `json:"-"`
	Verfication bool   `json:"verfication"`
}
