package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

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

func UpdateWSCryptoData(ws *websocket.Conn) {
	var cryptos []model.Crypto
	if err := Database.GetDB().Find(&cryptos).Error; err != nil {
		return
	}
	for i := range cryptos {
		exchangeData, err := GetExchangeRate(cryptos[i].Name)
		if err != nil {
			continue
		}

		cryptos[i].Name = exchangeData.From
		cryptos[i].Price = float64(exchangeData.AmountTo)
		if err := Database.GetDB().Save(&cryptos[i]).Error; err != nil {
		}
	}

	ws.WriteJSON(cryptos)
}