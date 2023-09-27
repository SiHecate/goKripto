package model

import "gorm.io/gorm"

type CryptoWallet struct {
	gorm.Model
	WalletID         int     `json:"wallet_address_id"`
	CryptoName       string  `json:"crypto_name"`
	CryptoTotalPrice float64 `json:"crypto_total_price"`
	Amount           float64 `json:"crypto_amount"`
	Wallet           Wallet  `gorm:"foreignKey:WalletID"`
}
