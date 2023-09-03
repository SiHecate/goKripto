package Database

import model "gokripto/Model"

func migrateTransaction() {
	GetDB().AutoMigrate(&model.Transaction{})
}
