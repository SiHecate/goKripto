package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func CryptoData(c *fiber.Ctx) error {
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
