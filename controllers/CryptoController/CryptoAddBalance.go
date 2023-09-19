package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func AddBalanceCrypto(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
	}

	addBalance := data["addBalance"].(float64)

	var availableBalance float64
	if err := Database.DB.Model(model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &availableBalance).Error; err != nil {
		return err
	}

	TotalBalance := addBalance + availableBalance

	if err := Database.DB.Model(model.Wallet{}).Where("user_id = ?", issuer).Update("Balance", TotalBalance).Error; err != nil {
		return err
	}

	TransactionBalance(c, issuer, addBalance, "Deposit", "Balance Adding")
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
