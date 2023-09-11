package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

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

	var WalletAddress string
	Database.GetDB().Model(model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress)

	var wallet model.Wallet
	Database.GetDB().Where("wallet_address = ?", WalletAddress).First(&wallet)

	return c.JSON(wallet)
}
