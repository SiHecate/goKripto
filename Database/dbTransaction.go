package Database

import model "gokripto/Model"

func migrateTransaction() {
	DB.AutoMigrate(&model.Transaction{})
}
