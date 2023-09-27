package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"user_name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
	Wallet   Wallet `json:"wallet" gorm:"foreignKey:UserID"`
}
