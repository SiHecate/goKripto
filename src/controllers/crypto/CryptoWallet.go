package controllers

import (
	model "gokripto/Model"
	"gokripto/database"
)

func CryptoWallet(CryptoID uint, CryptoName string, CryptoPrice float64, Amount float64, WalletAddress string, ProcessType string) {
	var existingCryptoWallet model.CryptoWallet
	result := database.DB.Where("wallet_address = ? AND crypto_name = ?", WalletAddress, CryptoName).First(&existingCryptoWallet)

	// Calculate the total price of the crypto being bought or sold.
	CryptoTotalPrice := CryptoPrice * Amount

	// Create a new crypto wallet entry with the provided data.
	newCryptoWallet := model.CryptoWallet{
		CryptoID:         CryptoID,         // Stable
		CryptoName:       CryptoName,       // Stable
		CryptoTotalPrice: CryptoTotalPrice, // Stable
		WalletAddress:    WalletAddress,    // Stable
		Amount:           Amount,           // Stable
	}

	// If the process type is "buy," update the existing crypto wallet entry or create a new one.
	if ProcessType == "buy" {
		if result.Error == nil {
			existingCryptoWallet.Amount += Amount
			existingCryptoWallet.CryptoTotalPrice += CryptoTotalPrice
			database.DB.Save(&existingCryptoWallet)
		} else {
			database.DB.Create(&newCryptoWallet)
		}
	} else if ProcessType == "sell" {
		// If the process type is "sell," update the existing crypto wallet entry or create a new one.
		if result.Error == nil {
			existingCryptoWallet.Amount -= Amount
			existingCryptoWallet.CryptoTotalPrice -= CryptoTotalPrice
			database.DB.Save(&existingCryptoWallet)
		} else {
			database.DB.Create(&newCryptoWallet)
		}
	}
}
