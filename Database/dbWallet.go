package Database

import model "gokripto/Model"

func migrateWallet() {
	GetDB().AutoMigrate(&model.Wallet{})
}
