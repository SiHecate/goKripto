package helpers

import (
	model "gokripto/Model"
	"gokripto/database"
)

func GetWalletAddress(issuer string) (string, error) {
	modelWallet := model.Wallet{}

	// Assuming you want to find the wallet address based on the issuer's user_id
	if err := database.DB.Debug().Where("user_id = ?", issuer).First(&modelWallet).Error; err != nil {
		return "", err
	}

	WalletAddress := modelWallet.WalletAddress
	return WalletAddress, nil
}
