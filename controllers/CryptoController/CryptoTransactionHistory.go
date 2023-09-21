package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

// List function for transaction of crypto
func TransactionListCrypto(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var TransactionCryptos []model.TransactionCrypto
	if err := Database.DB.Where("user_id = ?", issuer).Find(&TransactionCryptos).Error; err != nil {
		return err
	}

	return c.JSON(TransactionCryptos)
}
