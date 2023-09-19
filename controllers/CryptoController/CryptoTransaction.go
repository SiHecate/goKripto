package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TransactionCryptos(c *fiber.Ctx, UserID string, price float64, cryptoname string, amount float64, transactionType string) error {
	Transaction := model.Transaction{
		UserID:     UserID,
		Price:      float64(price),
		CryptoName: cryptoname,
		Amount:     float64(amount),
		Type:       transactionType,
		Date:       time.Now(),
	}

	if err := Database.DB.Create(&Transaction).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Transaction successful",
	})
}
