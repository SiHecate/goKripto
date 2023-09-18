package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func ListAllCryptos(c *fiber.Ctx) error {
	// UpdateCryptoData(c)
	var cryptos []model.Crypto
	Database.GetDB().Find(&cryptos)
	return c.JSON(cryptos)
}

func ListCryptoWallet(c *fiber.Ctx) error {
	// UpdateCryptoData(c)
	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var WalletAddress string
	if err := Database.GetDB().Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	var cryptoWallets []model.CryptoWallet
	if err := Database.GetDB().Model(&model.CryptoWallet{}).Where("wallet_address = ? AND amount > ? ", WalletAddress, 0).Find(&cryptoWallets).Error; err != nil {
		return err
	}

	return c.JSON(cryptoWallets)
}
