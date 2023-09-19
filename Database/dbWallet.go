package Database

import model "gokripto/Model"

func migrateWallet() {
	DB.AutoMigrate(&model.Wallet{})
}
