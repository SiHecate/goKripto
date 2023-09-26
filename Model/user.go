package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"user_name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
	Wallet   Wallet `json:"wallet" gorm:"foreignKey:ID"`
}
