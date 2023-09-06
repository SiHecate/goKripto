package model

import (
	"time"
)

type Transaction struct {
	ID         uint      `json:"transactionID" gorm:"unique"`
<<<<<<< HEAD
	User       string    `json:"userID"`
	CryptoName string    `json:"cryptoname"`
	Price      float64   `json:"price"`
	Amount     float64   `json:"quantity"`
	Type       string    `json:"type"`
=======
	Name       string    `json:"username"`
	Price      int       `json:"price"`
	CryptoID   uint      `json:"cryptoID"`
	CryptoName string    `json:"cryptoname"`
	Quantity   int       `json:"quantity"`
>>>>>>> origin/main
	Date       time.Time `json:"-"`
}
