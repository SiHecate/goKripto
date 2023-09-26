package model

type Wallet struct {
	ID            uint         `json:"id" gorm:"primaryKey"`
	WalletAddress string       `json:"wallet_address" gorm:"unique"`
	UserID        uint         `json:"user_id" gorm:"not null;index"`
	Balance       float64      `json:"balance"`
	CryptoWallet  CryptoWallet `json:"crypto_wallet" gorm:"foreignKey:WalletAddress"`
}
