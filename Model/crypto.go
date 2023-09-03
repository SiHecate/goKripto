package model

import (
	"time"
)

type Crypto struct {
	ID        uint      `json:"cryptoID" gorm:"primaryKey"`
	Name      string    `json:"cryptoname"`
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"-"`
}
