package model

import (
	"time"

	"gorm.io/gorm"
)

// CryptoAPIResponse model the response from a cryptocurrency API.
type CryptoAPIResponse struct {
	Data []CryptoData `json:"data"`
}

// CryptoData model data for a single cryptocurrency.
type CryptoData struct {
	ID                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
	Explorer          string `json:"explorer"`
}

// Crypto model a cryptocurrency.
type Crypto struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Symbol    string    `json:"crypto_symbol"`
	Name      string    `json:"crypto_name"`
	Price     float64   `json:"crypto_price"`
	UpdatedAt time.Time `json:"-"`
}

// User model user information.
type User struct {
	gorm.Model
	Name     string `json:"user_name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
	Wallet   Wallet `json:"wallet" gorm:"foreignKey:UserID"`
}

// Wallet model a balance wallet.
type Wallet struct {
	gorm.Model
	WalletAddress string  `json:"wallet_address" gorm:"unique"`
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

// TransactionBalance model a balance transaction.
type TransactionBalance struct {
	ID            uint      `json:"id" gorm:"unique"`
	UserID        string    `json:"user_id"`
	WalletAddress string    `json:"wallet_address"`
	BalanceAmount float64   `json:"balance_amount"`
	Type          string    `json:"type"`
	TypeInfo      string    `json:"type_info"`
	Date          time.Time `json:"date"`
}

// TransactionCrypto model a cryptocurrency transaction.
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
