package model

type Wallet struct {
	ID            uint     `json:"walletID" gorm:"primaryKey"`
	WalletAddress string   `json:"walletAddress" gorm:"unique"`
	UserID        uint     `json:"userID"` // Kullanıcının kimliği (User tablosuyla ilişkili)
	Balance       float32  `json:"balance"`
	Cryptos       []Crypto `json:"cryptos" gorm:"many2many:wallet_crypto;"`
}
