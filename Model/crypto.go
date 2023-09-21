package model

import (
	"time"
)

type Crypto struct {
	ID        uint      `json:"cryptoID" gorm:"primaryKey"`
	Symbol    string    `json:"cryptoSymbol"`
	Name      string    `json:"cryptoName"`
	Price     float64   `json:"cryptoPrice"`
	UpdatedAt time.Time `json:"-"`
}
