package controllers

import (
	model "gokripto/Model"
	"gokripto/database"

	"github.com/gofiber/fiber/v2"
)

func AddBalanceCrypto(c *fiber.Ctx) error {
	// Get the issuer (user ID) from the JWT token.
	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// Parsing
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
	}

	addBalance := data["addBalance"].(float64)

	var WalletAddress string
	if err := database.DB.Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	var availableBalance float64
	if err := database.DB.Model(model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &availableBalance).Error; err != nil {
		return err
	}

	// Total balance calculation
	TotalBalance := addBalance + availableBalance

	if err := database.DB.Model(model.Wallet{}).Where("user_id = ?", issuer).Update("Balance", TotalBalance).Error; err != nil {
		return err
	}

	//JSON response
	TransactionBalance(c, issuer, WalletAddress, addBalance, "Deposit", "Balance Adding")
	type addBalanceResponse struct {
		Issuer           string  `json:"issuer"`
		AvailableBalance float64 `json:"availableBalance"`
		TotalBalance     float64 `json:"addedBalance"`
	}

	response := addBalanceResponse{
		Issuer:           issuer,
		AvailableBalance: availableBalance,
		TotalBalance:     TotalBalance,
	}

	return c.JSON(response)
}
