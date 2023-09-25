package controllers

import (
	model "gokripto/Model"
	"gokripto/database"

	"github.com/gofiber/fiber/v2"
)

func AccountBalance(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	// Retrieve the wallet address associated with the user from the database.
	var WalletAddress string
	if err := database.DB.Model(model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	// Retrieve the wallet information (balance, etc.) based on the wallet address.
	var wallet model.Wallet
	if err := database.DB.Where("wallet_address = ?", WalletAddress).First(&wallet).Error; err != nil {
		return err
	}

	return c.JSON(wallet)
}
