package Database

import model "gokripto/Model"

func migrateTransactionBalance() {
	DB.AutoMigrate(&model.TransactionBalance{})
}
