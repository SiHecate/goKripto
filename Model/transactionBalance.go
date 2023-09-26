package model

import (
	"time"
)

type TransactionBalance struct {
	ID            uint      `json:"id" gorm:"unique"`
	UserID        string    `json:"user_id"`
	WalletAddress string    `json:"wallet_address"`
	BalanceAmount float64   `json:"balance_amount"`
	Type          string    `json:"type"`
	TypeInfo      string    `json:"type_info"`
	Date          time.Time `json:"date"`
}
