package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var cryptoNames = []string{"btc", "eth", "ltc"}

func AddCryptoData(c *fiber.Ctx) error {
	for _, cryptoName := range cryptoNames {
		exchangeData, err := GetExchangeRate(cryptoName)
		if err != nil {
			return err
		}

		crypto := model.Crypto{
			Name:  exchangeData.FromNetwork,
			Price: float64(exchangeData.AmountTo),
		}

		if err := Database.GetDB().Create(&crypto).Error; err != nil {
			return err
		}
	}

	return c.JSON("Crypto data added successfully")
}

func UpdateCryptoData(c *fiber.Ctx) error {
	var cryptos []model.Crypto
	if err := Database.GetDB().Find(&cryptos).Error; err != nil {
		return err
	}
	for i := range cryptos {
		exchangeData, err := GetExchangeRate(cryptoNames[i])
		if err != nil {
			return err
		}

		cryptos[i].Name = exchangeData.From
		cryptos[i].Price = float64(exchangeData.AmountTo)
		if err := Database.GetDB().Save(&cryptos[i]).Error; err != nil {
			return err
		}
	}

	return c.JSON(cryptos)
}

func ListAllCryptos(c *fiber.Ctx) error {
	var cryptos []model.Crypto
	Database.GetDB().Find(&cryptos)
	return c.JSON(cryptos)
}

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

	var userBalance float64
	Database.GetDB().Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &userBalance)
	var cryptoPrice float64
	Database.GetDB().Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice)
	totalCost := cryptoPrice * amountToBuy
	totalBalance := userBalance - totalCost
	Database.GetDB().Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance)

	TransactionCryptos(c, issuer, cryptoPrice, cryptoName, amountToBuy, "Buy")

	return c.JSON(fiber.Map{
		"totalCost":     totalCost,
		"cryptoName":    cryptoName,
		"amountToBuy":   amountToBuy,
		"iss":           issuer,
		"userbalance":   userBalance,
		"userbalanceAB": totalBalance,
	})
}

func SellCryptos(c *fiber.Ctx) error {
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
	totalGain := cryptoPrice * amountToSell
	totalBalance := userBalance + totalGain
	Database.GetDB().Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance)

	TransactionCryptos(c, issuer, cryptoPrice, cryptoName, amountToSell, "Sell")

	return c.JSON(fiber.Map{
		"totalGain":     totalGain,
		"cryptoName":    cryptoName,
		"amountToSell":  amountToSell,
		"iss":           issuer,
		"userbalance":   userBalance,
		"userbalanceAB": totalBalance,
	})
}

func TransactionCryptos(c *fiber.Ctx, UserID string, price float64, cryptoname string, amount float64, transactionType string) error {
	Transaction := model.Transaction{
		User:       UserID,
		Price:      float64(price),
		CryptoName: cryptoname,
		Amount:     float64(amount),
		Type:       transactionType,
	}

	Database.GetDB().Create(&Transaction)
	return c.JSON(fiber.Map{
		"message": "Transaction successful",
	})
}
