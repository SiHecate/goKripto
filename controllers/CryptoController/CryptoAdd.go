package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

var cryptoNames = []string{"btc", "eth", "ltc", "xmr", "bch", "bnb"}

func AddCryptoData(c *fiber.Ctx) error {
	for _, cryptoName := range cryptoNames {
		if cryptoExists(cryptoName) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Crypto already exists in the database.",
			})
		}

		if err := createCrypto(cryptoName); err != nil {
			return err
		}
	}

	return c.JSON(fiber.Map{
		"message": "Crypto data added successfully",
	})
}

func cryptoExists(name string) bool {
	var existingCrypto model.Crypto
	if err := Database.DB.Where("name = ?", name).First(&existingCrypto).Error; err == nil {
		return true
	}
	return false
}

func createCrypto(name string) error {
	exchangeData, err := GetExchangeRate(name)
	if err != nil {
		return err
	}

	crypto := model.Crypto{
		Name:  name,
		Price: float64(exchangeData.AmountTo),
	}

	if err := Database.DB.Create(&crypto).Error; err != nil {
		return err
	}

	return nil
}
