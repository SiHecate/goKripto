package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Create a new transaction balance entry.
func TransactionBalance(c *fiber.Ctx, UserID string, WalletAddress string, BalanceAmount float64, transactionType string, transactionInfo string) error {
	TransactionBalance := model.TransactionBalance{
		UserID:        UserID,
		WalletAddress: WalletAddress,
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
