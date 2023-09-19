package Database

import model "gokripto/Model"

func migrateTransactionCrypto() {
	DB.AutoMigrate(&model.TransactionCrypto{})
}
