package controllers

import (
	"encoding/json"
	"gokripto/Database"
	model "gokripto/Model"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func AddAllCryptoData(c *fiber.Ctx) error {
	apiURL := "https://api.coincap.io/v2/assets"

	response, err := http.Get(apiURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch crypto data from the API.",
			"error":   err.Error(),
		})
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "HTTP Error: " + strconv.Itoa(response.StatusCode),
		})
	}

	var cryptoAPIResponse model.CryptoAPIResponse
	err = json.NewDecoder(response.Body).Decode(&cryptoAPIResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse API response.",
			"error":   err.Error(),
		})
	}

	for _, crypto := range cryptoAPIResponse.Data {
		if err := createCrypto(crypto.Symbol, crypto.Name, crypto.PriceUsd); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to add or update crypto data in the database.",
				"error":   err.Error(),
			})
		}
	}

	return c.JSON(fiber.Map{
		"message": "All crypto data added or updated successfully",
	})
}

func createCrypto(symbol string, name string, price string) error {
	// Dönüşüm işlemi: API'den gelen price stringini bir float64'e çevirin
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return err
	}

	// Kayıt bulunamadıysa yeni kayıt ekle
	crypto := model.Crypto{
		Symbol: symbol,
		Name:   name,
		Price:  priceFloat,
	}
	if err := Database.DB.Create(&crypto).Error; err != nil {
		return err
	}

	return nil
}
