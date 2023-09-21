package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func SellCryptos(c *fiber.Ctx) error {

	// Extract the issuer (user ID) from the JWT token.
	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Parse the request JSON data.
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Extract cryptoName and amountToSell from the request data.
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
			"message": "Invalid cryptoAmount",
		})
	}

	// Fetch the user's wallet and balance from the database using GORM.
	var wallet model.Wallet
	var userBalance float64
	if err := Database.DB.Where("user_id = ?", issuer).First(&wallet).Error; err != nil {
		return err
	}
	wallet.Balance = userBalance

	// Fetch the price of the selected cryptocurrency from the database.
	var cryptoPrice float64
	Database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice)

	// Calculate the total profit from the sale and update the user's balance.
	totalProfit := cryptoPrice * amountToSell
	totalBalance := userBalance + totalProfit
	if err := Database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	// Fetch the ID of the selected cryptocurrency and the user's wallet address.
	var cryptoID uint
	if err := Database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID).Error; err != nil {
		return err
	}
	var WalletAddress string
	if err := Database.DB.Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	// Fetch the cryptocurrency amount from the crypto wallet and check for invalid values.
	var cryptocurrency float64
	Database.DB.Model(&model.CryptoWallet{}).Where("wallet_address = ? AND crypto_name = ?", WalletAddress, cryptoName).Pluck("amount", &cryptocurrency)
	if cryptocurrency < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid number",
		})
	}

	// Perform transactions and update the crypto wallet.
	TransactionBalance(c, issuer, WalletAddress, totalProfit, "Sales", "Crypto Sales")
	TransactionCryptos(c, issuer, WalletAddress, cryptoPrice, cryptoName, amountToSell, "Sell")
	CryptoWallet(cryptoID, cryptoName, cryptoPrice, amountToSell, WalletAddress, "sell")

	// Define a response structure.
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

	// Create a response object.
	response := SellCryptoResponse{
		TotalProfit:      totalProfit,
		CryptoName:       cryptoName,
		AmountToSell:     amountToSell,
		Issuer:           issuer,
		UserBalance:      userBalance,
		UserBalanceAfter: totalBalance,
		Currency:         cryptocurrency,
	}

	// Return a JSON response with the sale details.
	return c.Status(200).JSON(response)
}
