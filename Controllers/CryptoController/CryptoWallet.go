package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"
)

func CryptoWallet(CryptoID int, CryptoName string, CryptoPrice float64, Amount float64, WalletAddress string, ProcessType string) {
	var existingCryptoWallet model.CryptoWallet
	result := Database.GetDB().Where("wallet_address = ? AND crypto_name = ?", WalletAddress, CryptoName).First(&existingCryptoWallet)

	if ProcessType == "buy" {
		if result.Error != nil {
			CryptoTotalPrice := CryptoPrice * Amount
			CryptoWallet := model.CryptoWallet{
				CryptoID:         CryptoID,
				CryptoName:       CryptoName,
				CryptoTotalPrice: CryptoTotalPrice,
				WalletAddress:    WalletAddress,
				Amount:           Amount,
			}
			Database.GetDB().Create(&CryptoWallet)
		} else {
			existingCryptoWallet.Amount += Amount
			existingCryptoWallet.CryptoTotalPrice = CryptoPrice * existingCryptoWallet.Amount
			Database.GetDB().Save(&existingCryptoWallet)
		}
	} else if ProcessType == "sell" {
		if result.Error != nil {
			CryptoTotalPrice := CryptoPrice * Amount
			CryptoWallet := model.CryptoWallet{
				CryptoID:         CryptoID,
				CryptoName:       CryptoName,
				CryptoTotalPrice: CryptoTotalPrice,
				WalletAddress:    WalletAddress,
				Amount:           Amount,
			}
			Database.GetDB().Create(&CryptoWallet)
		} else {
			existingCryptoWallet.Amount -= Amount
			existingCryptoWallet.CryptoTotalPrice = CryptoPrice * existingCryptoWallet.Amount
			Database.GetDB().Save(&existingCryptoWallet)
		}

	}
}
