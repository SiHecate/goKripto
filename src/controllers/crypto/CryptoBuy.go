package controllers

import (
	model "gokripto/Model"
	"gokripto/database"

	"github.com/gofiber/fiber/v2"
)

const SecretKey = "secret"

// BuyCryptos handles the purchase of cryptocurrencies.
func BuyCryptos(c *fiber.Ctx) error {

	// Get the issuer (user ID) from the JWT token.
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

	// Extract cryptoName and amountToBuy from the request data.
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

	// Fetch the price of the selected cryptocurrency from the database.
	var cryptoPrice float64
	if err := database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice).Error; err != nil {
		return err
	}

	// Fetch the user's current balance from the wallet.
	var userBalance float64
	if err := database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &userBalance).Error; err != nil {
		return err
	}

	// Calculate the total cost of the purchase.
	totalCost := cryptoPrice * amountToBuy
	totalBalance := userBalance - totalCost

	// Check if the user has sufficient balance for the purchase.
	if totalCost > userBalance {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Insufficient balance",
			"balance": userBalance,
			"cost":    totalCost,
		})
	}

	// Update the user's wallet balance.
	if err := database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	// Fetch the cryptocurrency's ID and user's wallet address.
	var cryptoID uint
	if err := database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID).Error; err != nil {
		return err
	}
	var WalletAddress string
	if err := database.DB.Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	// Perform transactions and update the crypto wallet.
	TransactionBalance(c, issuer, WalletAddress, totalCost, "Purchase", "Crypto Purchase")
	TransactionCryptos(c, issuer, WalletAddress, cryptoPrice, cryptoName, amountToBuy, "Buy")
	CryptoWallet(cryptoID, cryptoName, cryptoPrice, amountToBuy, WalletAddress, "buy")

	// Define a response structure.
	type BuyCryptoResponse struct {
		TotalCost     float64 `json:"totalCost"`
		CryptoName    string  `json:"cryptoName"`
		AmountToBuy   float64 `json:"amountToBuy"`
		Issuer        string  `json:"issuer"`
		UserBalance   float64 `json:"userBalance"`
		UserBalanceAB float64 `json:"newUserBalance"`
	}

	// Create a response object.
	response := BuyCryptoResponse{
		TotalCost:     totalCost,
		CryptoName:    cryptoName,
		AmountToBuy:   amountToBuy,
		Issuer:        issuer,
		UserBalance:   userBalance,
		UserBalanceAB: totalBalance,
	}

	// Return a JSON response with the purchase details.
	return c.Status(200).JSON(response)
}
