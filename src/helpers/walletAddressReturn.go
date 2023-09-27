package helpers

import (
	model "gokripto/Model"
	"gokripto/database"
)

func GetWalletAddress(issuer string) (string, error) {
	modelWallet := model.Wallet{}

	if err := database.DB.Where("user_id = ?", issuer).First(&modelWallet).Error; err != nil {
		return "", err
	}

	WalletAddress := modelWallet.WalletAddress
	return WalletAddress, nil
}
