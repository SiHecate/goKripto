package model

type Wallet struct {
	ID            uint         `json:"id" gorm:"primaryKey"`
	WalletAddress string       `json:"wallet_address" gorm:"unique"`
	UserID        uint         `json:"user_id"`
	Balance       float64      `json:"balance"`
	CryptoWallet  CryptoWallet ` gorm:"foreignKey:WalletAddress" json:"crypto_wallet,omitempty"`
}
