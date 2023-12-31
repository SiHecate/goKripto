package model

import "gorm.io/gorm"

// TransactionBalance model a balance transaction.
type TransactionBalance struct {
	gorm.Model
	ID            uint    `json:"id" gorm:"unique"`
	UserID        string  `json:"user_id"`
	WalletAddress string  `json:"wallet_address"`
	BalanceAmount float64 `json:"balance_amount"`
	Type          string  `json:"type"`
	TypeInfo      string  `json:"type_info"`
}

// TransactionCrypto model a cryptocurrency transaction.
type TransactionCrypto struct {
	gorm.Model
	ID           uint    `json:"id" gorm:"unique"`
	UserID       string  `json:"user_id"`
	WalletAddres string  `json:"wallet_address"`
	CryptoName   string  `json:"crypto_name"`
	Price        float64 `json:"crypto_price"`
	Amount       float64 `json:"crypto_amount"`
	Type         string  `json:"transaction_type"`
}
