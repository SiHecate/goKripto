package model

import "gorm.io/gorm"

type Wallet struct {
	gorm.Model
	WalletAddress string  `json:"wallet_address" gorm:"unique"`
	UserID        uint    `json:"user_id" gorm:"not null;index"`
	Balance       float64 `json:"balance"`
}
