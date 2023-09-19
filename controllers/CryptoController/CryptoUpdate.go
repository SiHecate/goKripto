package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// UpdateCryptoData şu an için işlevsiz çünkü Websocket sürekli olarak değerleri günceleyecek.
func UpdateCryptoData(c *fiber.Ctx) error {
	var cryptos []model.Crypto
	if err := Database.DB.Find(&cryptos).Error; err != nil {
		return err
	}
	for i := range cryptos {
		exchangeData, err := GetExchangeRate(cryptoNames[i])
		if err != nil {
			return err
		}

		cryptos[i].Name = exchangeData.From
		cryptos[i].Price = float64(exchangeData.AmountTo)
		if err := Database.DB.Save(&cryptos[i]).Error; err != nil {
			return err
		}
	}

	return c.JSON(cryptos)
}

func UpdateWSCryptoData(ws *websocket.Conn) {
	var cryptos []model.Crypto
	if err := Database.DB.Find(&cryptos).Error; err != nil {
		return
	}
	for i := range cryptos {
		exchangeData, err := GetExchangeRate(cryptos[i].Name)
		if err != nil {
			continue
		}

		// Yeni verilerle mevcut kaydı güncelle ve "name" sütununda çakışma durumunda mevcut kaydı güncelle
		if err := Database.DB.Model(&cryptos[i]).Where("name = ?", cryptos[i].Name).Updates(model.Crypto{
			Name:  exchangeData.From,
			Price: float64(exchangeData.AmountTo),
		}).Error; err != nil {
		}
	}

	ws.WriteJSON(cryptos)
}
