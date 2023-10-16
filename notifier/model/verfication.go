package model

import (
	"gorm.io/gorm"
)

type Verfication struct {
	gorm.Model
	Email       string
	Verfication string
}
