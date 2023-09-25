package controllers

import (
	"encoding/json"
	model "gokripto/Model"
	"gokripto/database"
	"net/http"
	"strconv"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
)

func AddAllCryptoData(ws *websocket.Conn) error {
	apiURL := "https://api.coincap.io/v2/assets"

	response, err := http.Get(apiURL)
	if err != nil {
		errorMessage := fiber.Map{
			"message": "Failed to fetch crypto data from the API.",
			"error":   err.Error(),
		}
		sendJSONError(ws, errorMessage)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorMessage := fiber.Map{
			"message": "HTTP Error: " + strconv.Itoa(response.StatusCode),
		}
		sendJSONError(ws, errorMessage)
		return nil
	}

	var cryptoAPIResponse model.CryptoAPIResponse
	err = json.NewDecoder(response.Body).Decode(&cryptoAPIResponse)
	if err != nil {
		errorMessage := fiber.Map{
			"message": "Failed to parse API response.",
			"error":   err.Error(),
		}
		sendJSONError(ws, errorMessage)
		return err
	}

	for _, crypto := range cryptoAPIResponse.Data {
		if err := createOrUpdateCrypto(crypto.Symbol, crypto.Name, crypto.PriceUsd); err != nil {
			errorMessage := fiber.Map{
				"message": "Failed to add or update crypto data in the database.",
				"error":   err.Error(),
			}
			sendJSONError(ws, errorMessage)
			return err
		}
	}

	successMessage := fiber.Map{
		"message": "All crypto data added or updated successfully",
	}
	sendJSONResponse(ws, successMessage)

	return nil
}

func createOrUpdateCrypto(symbol string, name string, price string) error {
	// Float64 parser for string variable
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return err
	}

	var crypto model.Crypto
	result := database.DB.Where("symbol = ?", symbol).First(&crypto)
	// This switch case detects database condition
	switch {
	case result.Error != nil:
		// If crypto database is empty or not exists, create a new database.
		crypto := model.Crypto{
			Symbol: symbol,
			Name:   name,
			Price:  priceFloat,
		}
		if err := database.DB.Create(&crypto).Error; err != nil {
			return err
		}
	case result.Error != nil:
		// Other possible error
		return err

	default:
		// Crypto price update
		crypto.Price = priceFloat
		database.DB.Save(&crypto)
	}

	return nil
}

func sendJSONError(ws *websocket.Conn, errorMessage fiber.Map) {
	err := ws.WriteJSON(errorMessage)
	if err != nil {
		// Handle write error
	}
}

func sendJSONResponse(ws *websocket.Conn, responseMessage fiber.Map) {
	err := ws.WriteJSON(responseMessage)
	if err != nil {
		// Handle write error
	}
}
