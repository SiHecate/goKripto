package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func TransactionListBalance(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var TransactionBalance []model.TransactionBalance
	if err := Database.DB.Where("user_id = ?", issuer).Find(&TransactionBalance).Error; err != nil {
		return err
	}
	return c.JSON(TransactionBalance)
}
