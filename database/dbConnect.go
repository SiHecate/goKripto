package database

import (
	model "gokripto/Model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func dbConnect() {
	dsn := "host=postgres user=postgres password=393406 dbname=kriptoDB port=5432 sslmode=disable TimeZone=Europe/Istanbul"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Database error: " + err.Error())
	}
}

func MigrateTables() {
	DB.AutoMigrate(
		&model.User{},
		&model.Crypto{},
		&model.CryptoWallet{},
		&model.Wallet{},
		&model.TransactionBalance{},
		&model.TransactionCrypto{},
	)
}

func Connect() {
	dbConnect()
	MigrateTables()
}
