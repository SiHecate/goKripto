package Database

import (
	model "gokripto/Model"
)

func migrateCryptoWallet() {
	DB.AutoMigrate(&model.CryptoWallet{})
}
