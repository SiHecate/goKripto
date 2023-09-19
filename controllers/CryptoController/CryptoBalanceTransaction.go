package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TransactionBalance(c *fiber.Ctx, UserID string, BalanceAmount float64, transactionType string, transactionInfo string) error {
	TransactionBalance := model.TransactionBalance{
		UserID:        UserID,
		BalanceAmount: float64(BalanceAmount),
		Type:          transactionType,
		TypeInfo:      transactionInfo,
		Date:          time.Now(),
	}

	if err := Database.DB.Create(&TransactionBalance).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Transaction successful",
	})
}
