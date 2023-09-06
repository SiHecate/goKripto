package Database

import model "gokripto/Model"

func migrateCrypto() {
	GetDB().AutoMigrate(&model.Crypto{})
}
