package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"
)

func CryptoWallet(CryptoID int, CryptoName string, CryptoPrice float64, Amount float64, WalletAddress string, ProcessType string) {
	var existingCryptoWallet model.CryptoWallet
	result := Database.DB.Where("wallet_address = ? AND crypto_name = ?", WalletAddress, CryptoName).First(&existingCryptoWallet)

	CryptoTotalPrice := CryptoPrice * Amount
	newCryptoWallet := model.CryptoWallet{
		CryptoID:         CryptoID,
		CryptoName:       CryptoName,
		CryptoTotalPrice: CryptoTotalPrice,
		WalletAddress:    WalletAddress,
		Amount:           Amount,
	}

	if ProcessType == "buy" {
		if result.Error == nil {
			existingCryptoWallet.Amount += Amount
			existingCryptoWallet.CryptoTotalPrice += CryptoTotalPrice
			Database.DB.Save(&existingCryptoWallet)
		} else {
			Database.DB.Create(&newCryptoWallet)
		}
	} else if ProcessType == "sell" {
		if result.Error == nil {
			existingCryptoWallet.Amount -= Amount
			existingCryptoWallet.CryptoTotalPrice -= CryptoTotalPrice
			Database.DB.Save(&existingCryptoWallet)
		} else {
			Database.DB.Create(&newCryptoWallet)
		}
	}
}
