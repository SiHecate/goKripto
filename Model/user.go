package model

type User struct {
	Id            uint   `json:"userID" gorm:"unique"`
	Name          string `json:"userName"`
	Email         string `json:"email" gorm:"unique"`
	Password      []byte `json:"-"`
	WalletAddress string `json:"walletAddress" gorm:"unique"`
}
