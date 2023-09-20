package model

import (
	"time"
)

type Crypto struct {
	ID        uint      `json:"cryptoID" gorm:"primaryKey"`
	Symbol    string    `json:"cryptoSymbol"`
	Name      string    `json:"cryptoname"`
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"-"`
}
