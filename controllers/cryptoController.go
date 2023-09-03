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

		crypto := model.Crypto{
			Name:  exchangeData.FromNetwork,
			Price: float64(exchangeData.AmountTo),
		}

		if err := Database.GetDB().Create(&crypto).Error; err != nil {
			return err
		}
	}

	return c.JSON("Crypto data added successfully")
}

func UpdateCryptoData(c *fiber.Ctx) error {
	var cryptos []model.Crypto
	if err := Database.GetDB().Find(&cryptos).Error; err != nil {
		return err
	}
	for i := range cryptos {
		exchangeData, err := GetExchangeRate(cryptoNames[i])
		if err != nil {
			return err
		}

		cryptos[i].Name = exchangeData.From
		cryptos[i].Price = float64(exchangeData.AmountTo)
		if err := Database.GetDB().Save(&cryptos[i]).Error; err != nil {
			return err
		}
	}

	return c.JSON(cryptos)
}
