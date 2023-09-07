package model

type CryptoWallet struct {
	ID               uint    `json:"transactionID" gorm:"unique"`
	WalletAddress    string  `json:"walletAddress"`
	CryptoID         int     `json:"cryptoID"`
	CryptoName       string  `json:"cryptoname"`
	CryptoTotalPrice float64 `json:"cryptoTotalPrice"`
	Amount           float64 `json:"cryptoAmount"`
}
