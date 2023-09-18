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
	if err := Database.GetDB().Model(model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	var wallet model.Wallet
	if err := Database.GetDB().Where("wallet_address = ?", WalletAddress).First(&wallet).Error; err != nil {
		return err
	}

	return c.JSON(wallet)
}
