package Database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func dbConnect() {
	dsn := "host=localhost user=postgres password=393406 dbname=kriptoDB port=5432 sslmode=disable TimeZone=Europe/Istanbul"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Veritabanına bağlanırken bir hata oluştu: " + err.Error())
	}
}

func GetDB() *gorm.DB {
	return db
}

func MigrateTables() {
	migrateUser()
	migrateCrypto()
	migrateTransaction()
	migrateWallet()
	migrateCryptoWallet()
}

func Connect() {
	dbConnect()
	MigrateTables()
}
