package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func SellCryptos(c *fiber.Ctx) error {
	// UpdateCryptoData(c)

	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	cryptoName, ok := data["cryptoName"].(string)
	if !ok {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid cryptoName field",
		})
	}
	amountToSell, ok := data["amountToSell"].(float64)
	if !ok {
		c.Status(fiber.StatusBadGateway)
		return c.JSON(fiber.Map{
			"message": "Invlaid cryptoAmount",
		})
	}

	var userBalance float64
	if err := Database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &userBalance).Error; err != nil {
		return err
	}

	var cryptoPrice float64
	Database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice)

	totalProfit := cryptoPrice * amountToSell
	totalBalance := userBalance + totalProfit

	if err := Database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	var cryptoID int
	if err := Database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID).Error; err != nil {
		return err
	}

	var WalletAddress string
	if err := Database.DB.Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	CryptoWallet(cryptoID, cryptoName, cryptoPrice, amountToSell, WalletAddress, "sell")

	var cryptocurrency float64
	Database.DB.Model(&model.CryptoWallet{}).Where("wallet_address = ? AND crypto_name = ?", WalletAddress, cryptoName).Pluck("amount", &cryptocurrency)
	if cryptocurrency < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid number",
		})
	}

	TransactionCryptos(c, issuer, cryptoPrice, cryptoName, amountToSell, "Sell")
	type SellCryptoResponse struct {
		TotalProfit      float64 `json:"totalProfit"`
		CryptoName       string  `json:"cryptoName"`
		AmountToSell     float64 `json:"amountToSell"`
		Issuer           string  `json:"issuer"`
		UserBalance      float64 `json:"userBalance"`
		UserBalanceAfter float64 `json:"userBalanceAB"`
		CryptoID         uint    `json:"cryptoID"`
		Currency         float64
	}
	response := SellCryptoResponse{
		TotalProfit:      totalProfit,
		CryptoName:       cryptoName,
		AmountToSell:     amountToSell,
		Issuer:           issuer,
		UserBalance:      userBalance,
		UserBalanceAfter: totalBalance,
		Currency:         cryptocurrency,
	}

	return c.Status(200).JSON(response)

}
