package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

const SecretKey = "secret"

func BuyCryptos(c *fiber.Ctx) error {
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
	amountToBuy := data["amountToBuy"].(float64)

	var cryptoPrice float64
	Database.GetDB().Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice)

	var userBalance float64
	Database.GetDB().Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &userBalance)

	totalCost := cryptoPrice * amountToBuy
	totalBalance := userBalance - totalCost
	Database.GetDB().Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance)

	//Crypto Wallet
	var cryptoID int
	Database.GetDB().Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID)
	var WalletAddress string
	Database.GetDB().Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress)

	TransactionCryptos(c, issuer, cryptoPrice, cryptoName, amountToBuy, "Buy")

	CryptoWallet(cryptoID, cryptoName, cryptoPrice, amountToBuy, WalletAddress, "buy")
	//CryptoWallet(CryptoID int, CryptoName string, CryptoPrice float64, Amount float64, WalletAddress string) {

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
