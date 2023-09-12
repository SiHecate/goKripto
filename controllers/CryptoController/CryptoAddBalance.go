package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func AddBalanceCrypto(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
	}

	addBalance := data["addBalance"].(float64)

	var availableBalance float64
	Database.GetDB().Model(model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &availableBalance)

	TotalBalance := addBalance + availableBalance

	Database.GetDB().Model(model.Wallet{}).Where("user_id = ?", issuer).Update("Balance", TotalBalance)

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
