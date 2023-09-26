package model

type CryptoWallet struct {
	ID               uint    `json:"id" gorm:"primaryKey"`
	WalletAddress    string  `json:"wallet_address"`
	CryptoName       string  `json:"crypto_name"`
	CryptoTotalPrice float64 `json:"crypto_total_price"`
	Amount           float64 `json:"crypto_amount"`
}
