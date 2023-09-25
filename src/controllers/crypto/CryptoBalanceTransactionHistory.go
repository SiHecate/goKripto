package controllers

import (
	model "gokripto/Model"
	"gokripto/database"

	"github.com/gofiber/fiber/v2"
)

// List function for transaction of balance
func TransactionListBalance(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var TransactionBalance []model.TransactionBalance
	if err := database.DB.Where("user_id = ?", issuer).Find(&TransactionBalance).Error; err != nil {
		return err
	}

	return c.JSON(TransactionBalance)
}
