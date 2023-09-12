package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

const SecretKey = "secret"

func BuyCryptos(c *fiber.Ctx) error {
	UpdateCryptoData(c)

	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	cryptoName := data["cryptoName"].(string)
	amountToBuy := data["amountToBuy"].(float64)

	var cryptoPrice float64
	Database.GetDB().Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice)

	var userBalance float64
	Database.GetDB().Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &userBalance)

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

	Database.GetDB().Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance)

	//Crypto Wallet
	var cryptoID int
	Database.GetDB().Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID)
	var WalletAddress string
	Database.GetDB().Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress)

	TransactionCryptos(c, issuer, cryptoPrice, cryptoName, amountToBuy, "Buy")
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
