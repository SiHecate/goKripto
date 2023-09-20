package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

const SecretKey = "secret"

func BuyCryptos(c *fiber.Ctx) error {
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

	amountToBuy, ok := data["amountToBuy"].(float64)
	if !ok {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid amountToBuy field",
		})
	}

	var cryptoPrice float64
	if err := Database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice).Error; err != nil {
		return err
	}

	var userBalance float64
	if err := Database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &userBalance).Error; err != nil {
		return err
	}

	totalCost := cryptoPrice * amountToBuy
	totalBalance := userBalance - totalCost

	if totalCost > userBalance {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Insufficient balance",
			"balance": userBalance,
			"cost":    totalCost,
		})
	}

	if err := Database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	//Crypto Wallet
	var cryptoID int
	if err := Database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID).Error; err != nil {
		return err
	}
	var WalletAddress string
	if err := Database.DB.Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	TransactionBalance(c, issuer, totalCost, "Purchase", "Crypto Purchase")
	TransactionCryptos(c, issuer, WalletAddress, cryptoPrice, cryptoName, amountToBuy, "Buy")
	CryptoWallet(cryptoID, cryptoName, cryptoPrice, amountToBuy, WalletAddress, "buy")

	type BuyCryptoResponse struct {
		TotalCost     float64 `json:"totalCost"`
		CryptoName    string  `json:"cryptoName"`
		AmountToBuy   float64 `json:"amountToBuy"`
		Issuer        string  `json:"issuer"`
		UserBalance   float64 `json:"userBalance"`
		UserBalanceAB float64 `json:"userBalanceAB"`
		CryptoID      int     `json:"cryptoID"`
	}
	response := BuyCryptoResponse{
		TotalCost:     totalCost,
		CryptoName:    cryptoName,
		AmountToBuy:   amountToBuy,
		Issuer:        issuer,
		UserBalance:   userBalance,
		UserBalanceAB: totalBalance,
	}

	return c.Status(200).JSON(response)
}
