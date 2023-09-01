package Database

import model "gokripto/Model"

func dbUser() {
	GetDB().AutoMigrate(&model.User{})
}
