package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func SellCryptos(c *fiber.Ctx) error {
	UpdateCryptoData(c)
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	claims := token.Claims.(*jwt.StandardClaims)
	issuer := claims.Issuer

	cryptoName := data["cryptoName"].(string)
	amountToSell := data["amountToSell"].(float64)

	var userBalance float64
	Database.GetDB().Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &userBalance)
	var cryptoPrice float64
	Database.GetDB().Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice)
	totalProfit := cryptoPrice * amountToSell
	totalBalance := userBalance + totalProfit
	Database.GetDB().Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance)

	var cryptoID int
	Database.GetDB().Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID)
	var WalletAddress string
	Database.GetDB().Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress)

	CryptoWallet(cryptoID, cryptoName, cryptoPrice, amountToSell, WalletAddress, "sell")

	var cryptocurrency float64
	Database.GetDB().Model(&model.CryptoWallet{}).Where("wallet_address = ? AND crypto_name = ?", WalletAddress, cryptoName).Pluck("amount", &cryptocurrency)
	if cryptocurrency < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid number",
		})
	}

	TransactionCryptos(c, issuer, cryptoPrice, cryptoName, amountToSell, "Sell")
	type SellCryptoResponse struct {
		TotalProfit   float64 `json:"totalProfit"`
		CryptoName    string  `json:"cryptoName"`
		AmountToSell  float64 `json:"amountToSell"`
		Issuer        string  `json:"issuer"`
		UserBalance   float64 `json:"userBalance"`
		UserBalanceAB float64 `json:"userBalanceAB"`
		CryptoID      uint    `json:"cryptoID"`
		Currency      float64
	}
	response := SellCryptoResponse{
		TotalProfit:   totalProfit,
		CryptoName:    cryptoName,
		AmountToSell:  amountToSell,
		Issuer:        issuer,
		UserBalance:   userBalance,
		UserBalanceAB: totalBalance,
		Currency:      cryptocurrency,
	}

	return c.Status(200).JSON(response)

}
