package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func AddCryptoData(c *fiber.Ctx) error {
	exchangeData, err := GetExchangeRate()
	if err != nil {
		return err
	}

	crypto := model.Crypto{
		Name:  exchangeData.FromNetwork,
		Price: int(exchangeData.AmountTo),
	}

	if err := Database.GetDB().Create(&crypto).Error; err != nil {
		return err
	}

	return c.JSON(crypto)
}

func UpdateCryptoData(c *fiber.Ctx) error {
	// Güncellenecek tüm paralar
	var cryptos []model.Crypto
	if err := Database.GetDB().Find(&cryptos).Error; err != nil {
		return err
	}
	for _, crypto := range cryptos {
		// Güncel veriyi çektirme
		exchangeData, err := GetExchangeRate()
		if err != nil {
			return err
		}

		// Güncellenecek değerler
		crypto.Name = exchangeData.FromNetwork
		crypto.Price = int(exchangeData.AmountTo)
		if err := Database.GetDB().Save(&crypto).Error; err != nil {
			return err
		}
	}

	return c.JSON(cryptos)
}
