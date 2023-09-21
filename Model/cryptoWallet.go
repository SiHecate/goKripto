package model

type CryptoWallet struct {
	ID               uint    `json:"cryptoWalletID" gorm:"unique"`
	WalletAddress    string  `json:"walletAddress"`
	CryptoID         uint    `json:"cryptoID"`
	CryptoName       string  `json:"cryptoName"`
	CryptoTotalPrice float64 `json:"cryptoTotalPrice"`
	Amount           float64 `json:"cryptoAmount"`
}
