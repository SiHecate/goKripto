package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

var cryptoNames = []string{"btc", "eth", "ltc"}

func AddCryptoData(c *fiber.Ctx) error {
	for _, cryptoName := range cryptoNames {
		exchangeData, err := GetExchangeRate(cryptoName)
		if err != nil {
			return err
		}

		var existingCrypto model.Crypto
		if err := Database.DB.Where("name = ?", exchangeData.FromNetwork).First(&existingCrypto).Error; err == nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Crypto already exists in the database.",
			})
		}

		crypto := model.Crypto{
			Name:  exchangeData.FromNetwork,
			Price: float64(exchangeData.AmountTo),
		}

		if err := Database.DB.Create(&crypto).Error; err != nil {
			return err
		}
	}

	return c.JSON(fiber.Map{
		"message": "Crypto data added successfully",
	})
}
