package database

import (
	"fmt"
	model "gokripto/Model"

	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	DBConnection()
	MigrateTables()
}

func DBConnection() {
	dsn := "host=postgres user=postgres password=393406 dbname=kriptoDB port=5432 sslmode=disable TimeZone=Europe/Istanbul"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Database error: " + err.Error())
	}
	fmt.Println(DB)
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
	DB.Migrator().DropTable(
		&model.User{},
		&model.Wallet{},
		&model.CryptoWallet{},
		&model.Crypto{},
		&model.TransactionBalance{},
		&model.TransactionCrypto{},
	)
}

// type MockDBConnection struct {
// 	mock.Mock
// }

// func (m *MockDBConnection) DBConnection() *gorm.DB {
// 	args := m.Called()
// 	return args.Get(0).(*gorm.DB)
// }

type MockMigrateTables struct {
	mock.Mock
}

func (m *MockMigrateTables) MigrateTables() {
	m.Called()
}

type MockDownTables struct {
	mock.Mock
}

func (m *MockDownTables) DownTables() {
	m.Called()
}
