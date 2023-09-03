package model

import (
	"time"
)

type Crypto struct {
	ID        uint      `json:"cryptoID" gorm:"primaryKey"`
	Name      string    `json:"cryptoname"`
	Price     int       `json:"price"`
	UpdatedAt time.Time `json:"-"`
}
