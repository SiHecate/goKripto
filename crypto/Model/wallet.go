package model

import "gorm.io/gorm"

// Wallet model a balance wallet.
type Wallet struct {
	gorm.Model
	WalletAddress string  `json:"wallet_address" gorm:"not null;index"`
	UserID        uint    `json:"user_id" gorm:"not null;index"`
	Balance       float64 `json:"balance"`
}

// CryptoWallet model a cryptocurrency wallet.
type CryptoWallet struct {
	gorm.Model
	WalletID         int     `json:"wallet_address_id"`
	CryptoName       string  `json:"crypto_name"`
	CryptoTotalPrice float64 `json:"crypto_total_price"`
	Amount           float64 `json:"crypto_amount"`
	Wallet           Wallet  `gorm:"foreignKey:WalletID" json:"wallet"`
}
