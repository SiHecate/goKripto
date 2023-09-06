package Database

import model "gokripto/Model"

func migrateUser() {
	GetDB().AutoMigrate(&model.User{})
}
