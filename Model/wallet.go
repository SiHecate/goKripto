package model

type Wallet struct {
	ID            uint    `json:"walletID" gorm:"primaryKey"`
	WalletAddress string  `json:"walletAddress" gorm:"unique"`
	UserID        uint    `json:"userID"` // Kullanıcının kimliği (User tablosuyla ilişkili)
	Balance       float64 `json:"balance"`
}
