package controllers

import (
	model "gokripto/Model"
	"gokripto/database"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TransactionCryptos(c *fiber.Ctx, UserID string, WalletAddres string, price float64, cryptoname string, amount float64, transactionType string) error {
	TransactionCryptos := model.TransactionCrypto{
		UserID:       UserID,
		WalletAddres: WalletAddres,
		Price:        float64(price),
		CryptoName:   cryptoname,
		Amount:       float64(amount),
		Type:         transactionType,
		Date:         time.Now(),
	}

	// Create crypto transaction database
	if err := database.DB.Create(&TransactionCryptos).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Transacton successful",
	})
}
