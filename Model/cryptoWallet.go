package model

type CryptoWallet struct {
	ID               uint    `json:"cryptoID" gorm:"primaryKey"`
	cryptoName       string  `json:"cryptowalletCryptoName"`
	cryptoTotalPrice float64 `json:"cryptowalletCryptoTotalPrice"`
	// Her crypto için ayrı bir yer tutulması gerekiyor
}
