package model

type Wallet struct {
	ID            uint    `json:"walletID" gorm:"primaryKey"`
	WalletAddress string  `json:"walletAddress" gorm:"unique"`
<<<<<<< HEAD
	UserID        uint    `json:"userID"`
=======
	UserID        uint    `json:"userID"` // Kullanıcının kimliği (User tablosuyla ilişkili)
>>>>>>> origin/main
	Balance       float64 `json:"balance"`
}
