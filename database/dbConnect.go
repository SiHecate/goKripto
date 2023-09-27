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
		&model.Wallet{},
		&model.CryptoWallet{},
		&model.Crypto{},
		&model.TransactionBalance{},
		&model.TransactionCrypto{},
	)

}

func DownTables() {
	DB.Migrator().DropTable(&model.User{})
	DB.Migrator().DropTable(&model.Wallet{})
	DB.Migrator().DropTable(&model.Crypto{})
	DB.Migrator().DropTable(&model.TransactionBalance{})
	DB.Migrator().DropTable(&model.TransactionCrypto{})
	DB.Migrator().DropTable(&model.CryptoWallet{})
}

func Connect() {
	dbConnect()
	MigrateTables()
}
