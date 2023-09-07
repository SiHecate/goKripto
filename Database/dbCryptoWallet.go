package Database

import model "gokripto/Model"

func migrateCryptoWallet() {
	GetDB().AutoMigrate(&model.CryptoWallet{})
}
