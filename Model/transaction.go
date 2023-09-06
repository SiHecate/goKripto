package model

import (
	"time"
)

type Transaction struct {
	ID         uint      `json:"transactionID" gorm:"unique"`
	User       string    `json:"userID"`
	CryptoName string    `json:"cryptoname"`
	Price      float64   `json:"price"`
	Amount     float64   `json:"quantity"`
	Type       string    `json:"type"`
	Date       time.Time `json:"-"`
}
