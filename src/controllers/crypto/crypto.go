package controllers

import (
	"encoding/json"
	model "gokripto/Model"
	"gokripto/database"
	"net/http"
	"strconv"
	"time"

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

func ListAllCryptos(c *fiber.Ctx) error {
	var cryptos []model.Crypto
	if err := database.DB.Find(&cryptos).Error; err != nil {
		return err
	}
	return c.JSON(cryptos)
}

func ListCryptoWallet(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		return err
	}

	var WalletAddress string
	if err := database.DB.Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	var cryptoWallets []model.CryptoWallet
	if err := database.DB.Model(&model.CryptoWallet{}).Where("wallet_address = ? AND amount > ? ", WalletAddress, 0).Find(&cryptoWallets).Error; err != nil {
		return err
	}

	return c.JSON(cryptoWallets)
}

func AccountBalance(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		return err
	}
	// Retrieve the wallet address associated with the user from the database.
	var WalletAddress string
	if err := database.DB.Model(model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	// Retrieve the wallet information (balance, etc.) based on the wallet address.
	var wallet model.Wallet
	if err := database.DB.Where("wallet_address = ?", WalletAddress).First(&wallet).Error; err != nil {
		return err
	}

	return c.JSON(wallet)
}

func AddBalanceCrypto(c *fiber.Ctx) error {
	// Get the issuer (user ID) from the JWT token.
	issuer, err := GetToken(c)
	if err != nil {
		return err
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

const SecretKey = "secret"

// BuyCryptos handles the purchase of cryptocurrencies.
func BuyCryptos(c *fiber.Ctx) error {

	// Get the issuer (user ID) from the JWT token.
	issuer, err := GetToken(c)
	if err != nil {
		return err
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

func SellCryptos(c *fiber.Ctx) error {

	// Extract the issuer (user ID) from the JWT token.
	issuer, err := GetToken(c)
	if err != nil {
		return err
	}

	// Parse the request JSON data.
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	// Extract cryptoName and amountToSell from the request data.
	cryptoName, ok := data["cryptoName"].(string)
	if !ok {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid cryptoName field",
		})
	}
	amountToSell, ok := data["amountToSell"].(float64)
	if !ok {
		c.Status(fiber.StatusBadGateway)
		return c.JSON(fiber.Map{
			"message": "Invalid cryptoAmount",
		})
	}

	// Fetch the user's wallet and balance from the database using GORM.
	var wallet model.Wallet
	var userBalance float64
	if err := database.DB.Where("user_id = ?", issuer).First(&wallet).Error; err != nil {
		return err
	}
	wallet.Balance = userBalance

	// Fetch the price of the selected cryptocurrency from the database.
	var cryptoPrice float64
	database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice)

	// Calculate the total profit from the sale and update the user's balance.
	totalProfit := cryptoPrice * amountToSell
	totalBalance := userBalance + totalProfit
	if err := database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	// Fetch the ID of the selected cryptocurrency and the user's wallet address.
	var cryptoID uint
	if err := database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID).Error; err != nil {
		return err
	}
	var WalletAddress string
	if err := database.DB.Model(&model.User{}).Where("id = ?", issuer).Pluck("wallet_address", &WalletAddress).Error; err != nil {
		return err
	}

	// Fetch the cryptocurrency amount from the crypto wallet and check for invalid values.
	var cryptocurrency float64
	database.DB.Model(&model.CryptoWallet{}).Where("wallet_address = ? AND crypto_name = ?", WalletAddress, cryptoName).Pluck("amount", &cryptocurrency)
	if cryptocurrency < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid number",
		})
	}

	// Perform transactions and update the crypto wallet.
	TransactionBalance(c, issuer, WalletAddress, totalProfit, "Sales", "Crypto Sales")
	TransactionCryptos(c, issuer, WalletAddress, cryptoPrice, cryptoName, amountToSell, "Sell")
	CryptoWallet(cryptoID, cryptoName, cryptoPrice, amountToSell, WalletAddress, "sell")

	// Define a response structure.
	type SellCryptoResponse struct {
		TotalProfit      float64 `json:"totalProfit"`
		CryptoName       string  `json:"cryptoName"`
		AmountToSell     float64 `json:"amountToSell"`
		Issuer           string  `json:"issuer"`
		UserBalance      float64 `json:"userBalance"`
		UserBalanceAfter float64 `json:"userBalanceAB"`
		CryptoID         uint    `json:"cryptoID"`
		Currency         float64
	}

	// Create a response object.
	response := SellCryptoResponse{
		TotalProfit:      totalProfit,
		CryptoName:       cryptoName,
		AmountToSell:     amountToSell,
		Issuer:           issuer,
		UserBalance:      userBalance,
		UserBalanceAfter: totalBalance,
		Currency:         cryptocurrency,
	}

	// Return a JSON response with the sale details.
	return c.Status(200).JSON(response)
}

// Create a new transaction balance entry.
func TransactionBalance(c *fiber.Ctx, UserID string, WalletAddress string, BalanceAmount float64, transactionType string, transactionInfo string) error {
	TransactionBalance := model.TransactionBalance{
		UserID:        UserID,
		WalletAddress: WalletAddress,
		BalanceAmount: float64(BalanceAmount),
		Type:          transactionType,
		TypeInfo:      transactionInfo,
		Date:          time.Now(),
	}

	if err := database.DB.Create(&TransactionBalance).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Transaction successful",
	})
}

// List function for transaction of balance
func TransactionListBalance(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		return err
	}

	var TransactionBalance []model.TransactionBalance
	if err := database.DB.Where("user_id = ?", issuer).Find(&TransactionBalance).Error; err != nil {
		return err
	}

	return c.JSON(TransactionBalance)
}

func TransactionCryptos(c *fiber.Ctx, UserID string, WalletAddres string, price float64, cryptoname string, amount float64, transactionType string) error {
	TransactionCryptos := model.TransactionCrypto{
		UserID:       UserID,
		WalletAddres: WalletAddres,
		Price:        float64(price),
		CryptoName:   cryptoname,
		Amount:       float64(amount),
		Type:         transactionType,
		Date:         time.Now(),
	}

	// Create crypto transaction database
	if err := database.DB.Create(&TransactionCryptos).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Transacton successful",
	})
}

func TransactionListCrypto(c *fiber.Ctx) error {
	issuer, err := GetToken(c)
	if err != nil {
		return err
	}

	var TransactionCryptos []model.TransactionCrypto
	if err := database.DB.Where("user_id = ?", issuer).Find(&TransactionCryptos).Error; err != nil {
		return err
	}

	return c.JSON(TransactionCryptos)
}

func CryptoWallet(CryptoID uint, CryptoName string, CryptoPrice float64, Amount float64, WalletAddress string, ProcessType string) {
	var existingCryptoWallet model.CryptoWallet
	result := database.DB.Where("wallet_address = ? AND crypto_name = ?", WalletAddress, CryptoName).First(&existingCryptoWallet)

	// Calculate the total price of the crypto being bought or sold.
	CryptoTotalPrice := CryptoPrice * Amount

	// Create a new crypto wallet entry with the provided data.
	newCryptoWallet := model.CryptoWallet{
		CryptoID:         CryptoID,         // Stable
		CryptoName:       CryptoName,       // Stable
		CryptoTotalPrice: CryptoTotalPrice, // Stable
		WalletAddress:    WalletAddress,    // Stable
		Amount:           Amount,           // Stable
	}

	// If the process type is "buy," update the existing crypto wallet entry or create a new one.
	if ProcessType == "buy" {
		if result.Error == nil {
			existingCryptoWallet.Amount += Amount
			existingCryptoWallet.CryptoTotalPrice += CryptoTotalPrice
			database.DB.Save(&existingCryptoWallet)
		} else {
			database.DB.Create(&newCryptoWallet)
		}
	} else if ProcessType == "sell" {
		// If the process type is "sell," update the existing crypto wallet entry or create a new one.
		if result.Error == nil {
			existingCryptoWallet.Amount -= Amount
			existingCryptoWallet.CryptoTotalPrice -= CryptoTotalPrice
			database.DB.Save(&existingCryptoWallet)
		} else {
			database.DB.Create(&newCryptoWallet)
		}
	}
}
