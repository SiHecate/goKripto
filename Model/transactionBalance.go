package model

import (
	"time"
)

type TransactionBalance struct {
	ID            uint      `json:"transactionID" gorm:"unique"`
	UserID        string    `json:"userID"`
	WalletAddress string    `json:"walletAddress"`
	BalanceAmount float64   `json:"balanceAmount"`
	Type          string    `json:"type"`
	TypeInfo      string    `json:"typeInfo"`
	Date          time.Time `json:"date"`
}
