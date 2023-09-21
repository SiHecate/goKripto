package model

import (
	"time"
)

type TransactionCrypto struct {
	ID           uint      `json:"transactionID" gorm:"unique"`
	UserID       string    `json:"userID"`
	WalletAddres string    `json:"walletAddress"`
	CryptoName   string    `json:"cryptoName"`
	Price        float64   `json:"cryptoPrice"`
	Amount       float64   `json:"cryptoAmount"`
	Type         string    `json:"transactionType"`
	Date         time.Time `json:"transactionDate"`
}
