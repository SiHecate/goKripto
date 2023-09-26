package model

import (
	"time"
)

type TransactionCrypto struct {
	ID           uint      `json:"id" gorm:"unique"`
	UserID       string    `json:"user_id"`
	WalletAddres string    `json:"wallet_address"`
	CryptoName   string    `json:"crypto_name"`
	Price        float64   `json:"crypto_price"`
	Amount       float64   `json:"crypto_amount"`
	Type         string    `json:"transaction_type"`
	Date         time.Time `json:"transaction_date"`
}
