package model

import (
	"time"
)

type Transaction struct {
	ID         uint      `json:"transactionID" gorm:"unique"`
	Name       string    `json:"username"`
	Price      int       `json:"price"`
	CryptoID   uint      `json:"cryptoID"`
	CryptoName string    `json:"cryptoname"`
	Quantity   int       `json:"quantity"`
	Date       time.Time `json:"-"`
}
