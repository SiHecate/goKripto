package Database

import model "gokripto/Model"

func migrateCrypto() {
	DB.AutoMigrate(&model.Crypto{})
}
