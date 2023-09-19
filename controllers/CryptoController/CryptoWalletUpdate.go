package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/contrib/websocket"
)

func UpdateWalletCryptos(ws *websocket.Conn) {
	var cryptoWallet []model.CryptoWallet
	if err := Database.DB.Find(&cryptoWallet).Error; err != nil {
		// Hata işleme kodu buraya gelebilir.
		return
	}

	for i := range cryptoWallet {
		exchangeData, err := GetExchangeRate(cryptoWallet[i].CryptoName)
		if err != nil {
			continue
		}

		var totalAmount float64
		if err := Database.DB.Model(&model.CryptoWallet{}).Where("crypto_name = ?", cryptoWallet[i].CryptoName).Pluck("amount", &totalAmount).Error; err != nil {
			continue
		}

		price := float64(exchangeData.AmountTo)
		newPrice := price * totalAmount

		if err := Database.DB.Model(&model.CryptoWallet{}).Where("crypto_name = ?", cryptoWallet[i].CryptoName).Updates(model.CryptoWallet{
			CryptoName:       exchangeData.From,
			CryptoTotalPrice: newPrice,
		}).Error; err != nil {
			continue
		}
	}
}
