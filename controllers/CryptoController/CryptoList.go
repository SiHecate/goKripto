package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
)

func ListAllCryptos(c *fiber.Ctx) error {
	UpdateCryptoData(c)
	var cryptos []model.Crypto
	Database.GetDB().Find(&cryptos)
	return c.JSON(cryptos)
}
