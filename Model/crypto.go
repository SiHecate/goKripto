package model

import (
	"time"
)

type Crypto struct {
	ID        uint      `json:"cryptoID" gorm:"primaryKey"`
	Name      string    `json:"cryptoname" gorm:"unique;not null"`
	Price     float64   `json:"price"`
	UpdatedAt time.Time `json:"-"`
}
