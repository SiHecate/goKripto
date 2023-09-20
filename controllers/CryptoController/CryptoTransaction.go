package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"
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

	if err := Database.DB.Create(&TransactionCryptos).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Transacton successful",
	})
}
