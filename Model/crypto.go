package model

import (
	"time"
)

type Crypto struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Symbol    string    `json:"crypto_symbol"`
	Name      string    `json:"crypto_name"`
	Price     float64   `json:"crypto_price"`
	UpdatedAt time.Time `json:"-"`
}
