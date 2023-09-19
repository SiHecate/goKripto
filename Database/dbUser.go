package Database

import model "gokripto/Model"

func migrateUser() {
	DB.AutoMigrate(&model.User{})
}
