package database

import (
	model "cryptoApp/Model"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Conn *gorm.DB

func Connect() {
	DBConnection()
	MigrateTables()

	fmt.Println("Database connection success!")
}

func DBConnection() {
	dsn := "host=postgres user=postgres password=393406 dbname=kriptoDB port=5432 sslmode=disable TimeZone=Europe/Istanbul"

	var err error
	Conn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Database error: " + err.Error())
	}
}

func MigrateTables() {
	Conn.AutoMigrate(
		&model.User{},
		&model.Wallet{},
		&model.CryptoWallet{},
		&model.Crypto{},
		&model.TransactionBalance{},
		&model.TransactionCrypto{},
		&model.Verfication{},
	)

}

func DownTables() {
	Conn.Migrator().DropTable(
		&model.User{},
		&model.Wallet{},
		&model.CryptoWallet{},
		&model.Crypto{},
		&model.TransactionBalance{},
		&model.TransactionCrypto{},
		&model.Verfication{},
	)
}
