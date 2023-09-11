package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func TransactionCryptos(c *fiber.Ctx, UserID string, price float64, cryptoname string, amount float64, transactionType string) error {
	Transaction := model.Transaction{
		User:       UserID,
		Price:      float64(price),
		CryptoName: cryptoname,
		Amount:     float64(amount),
		Type:       transactionType,
	}

	Database.GetDB().Create(&Transaction)
	return c.JSON(fiber.Map{
		"message": "Transaction successful",
	})
}
