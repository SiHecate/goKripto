package controllers

import (
	"encoding/json"
	"fmt"
	model "gokripto/Model"
	"gokripto/database"
	helpers "gokripto/pkg/helpers"
	"net/http"
	"strconv"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return err
	}

	var crypto model.Crypto
	result := database.DB.Where("symbol = ?", symbol).First(&crypto)
	switch {
	case result.Error != nil:
		crypto := model.Crypto{
			Symbol: symbol,
			Name:   name,
			Price:  priceFloat,
		}
		if err := database.DB.Create(&crypto).Error; err != nil {
			return err
		}
	case result.Error != nil:
		return err

	default:
		crypto.Price = priceFloat
		database.DB.Save(&crypto)
	}

	return nil
}

func sendJSONError(ws *websocket.Conn, errorMessage fiber.Map) {
	err := ws.WriteJSON(errorMessage)
	if err != nil {
	}
}

func sendJSONResponse(ws *websocket.Conn, responseMessage fiber.Map) {
	err := ws.WriteJSON(responseMessage)
	if err != nil {
	}
}

func ListAllCryptos(c *fiber.Ctx) error {
	var cryptos []model.Crypto
	if err := database.DB.Find(&cryptos).Error; err != nil {
		return err
	}

	type ListAllCrypto struct {
		Symbol string  `json:"crypto_symbol"`
		Name   string  `json:"crypto_name"`
		Price  float64 `json:"crypto_price"`
	}

	var response []ListAllCrypto
	for _, crypto := range cryptos {
		response = append(response, ListAllCrypto{
			Symbol: crypto.Symbol,
			Name:   crypto.Name,
			Price:  crypto.Price,
		})
	}
	return c.JSON(response)
}

func AccountBalance(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuerFromContext(c)
	if err != nil {
		return err
	}

	WalletAddress, err := helpers.GetWalletAddress(issuer)
	if err != nil {
		return err
	}

	var wallet model.Wallet
	if err := database.DB.Where("wallet_address = ?", WalletAddress).First(&wallet).Error; err != nil {
		return err
	}

	var user model.User
	if err := database.DB.Where("id = ?", issuer).Preload("Wallet").First(&user).Error; err != nil {
		return err
	}

	type WalletResponse struct {
		WalletAddress string  `json:"wallet_address"`
		Username      string  `json:"username"`
		Balance       float64 `json:"balance"`
	}

	response := WalletResponse{
		WalletAddress: wallet.WalletAddress,
		Username:      user.Name,
		Balance:       wallet.Balance,
	}

	return c.JSON(response)
}

func AddBalanceCrypto(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuerFromContext(c)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
	}

	addBalance := data["addBalance"].(float64)

	WalletAddress, err := helpers.GetWalletAddress(issuer)
	if err != nil {
		return err
	}

	var availableBalance float64
	if err := database.DB.Model(model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &availableBalance).Error; err != nil {
		return err
	}

	TotalBalance := addBalance + availableBalance

	if err := database.DB.Model(model.Wallet{}).Where("user_id = ?", issuer).Update("Balance", TotalBalance).Error; err != nil {
		return err
	}

	TransactionBalance(c, issuer, WalletAddress, addBalance, "Deposit", "Balance Adding")
	type addBalanceResponse struct {
		Issuer           string  `json:"issuer"`
		AvailableBalance float64 `json:"available_balance"`
		TotalBalance     float64 `json:"added_balance"`
	}

	response := addBalanceResponse{
		Issuer:           issuer,
		AvailableBalance: availableBalance,
		TotalBalance:     TotalBalance,
	}

	return c.JSON(response)
}

const SecretKey = "secret"

func BuyCryptos(c *fiber.Ctx) error {

	issuer, err := helpers.GetIssuerFromContext(c)
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

	var cryptoPrice float64
	if err := database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice).Error; err != nil {
		return err
	}
	fmt.Printf("Kullanıcı Kimliği (User ID): '%s'\n", issuer)

	var userBalance float64
	if err := database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &userBalance).Error; err != nil {
		return err
	}

	totalCost := cryptoPrice * amountToBuy
	totalBalance := userBalance - totalCost

	if totalCost > userBalance {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Insufficient balance",
			"balance": userBalance,
			"cost":    totalCost,
		})
	}

	if err := database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	var cryptoID uint
	if err := database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID).Error; err != nil {
		return err
	}

	modelWallet := model.Wallet{}
	database.DB.Debug().Where("user_id = ?", issuer).Find(&modelWallet)

	WalletAddress, err := helpers.GetWalletAddress(issuer)
	if err != nil {
		return err
	}

	CryptoWallet(issuer, cryptoName, cryptoPrice, amountToBuy, "buy")
	TransactionBalance(c, issuer, WalletAddress, totalCost, "Purchase", "Crypto Purchase")
	TransactionCryptos(c, issuer, WalletAddress, cryptoPrice, cryptoName, amountToBuy, "Buy")

	type BuyCryptoResponse struct {
		TotalCost     float64 `json:"totalCost"`
		CryptoName    string  `json:"cryptoName"`
		AmountToBuy   float64 `json:"amountToBuy"`
		Issuer        string  `json:"issuer"`
		UserBalance   float64 `json:"userBalance"`
		UserBalanceAB float64 `json:"newUserBalance"`
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

func SellCryptos(c *fiber.Ctx) error {

	issuer, err := helpers.GetIssuerFromContext(c)
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

	var userBalance float64
	if err := database.DB.Model(&model.Wallet{}).Where("id = ?", issuer).Pluck("balance", &userBalance).Error; err != nil {
		return err
	}

	var cryptoPrice float64
	database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice)

	totalProfit := cryptoPrice * amountToSell
	totalBalance := userBalance + totalProfit
	if err := database.DB.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	var cryptoID uint
	if err := database.DB.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID).Error; err != nil {
		return err
	}

	modelWallet := model.Wallet{}
	database.DB.Where("user_id = ?", issuer).Find(&modelWallet)

	WalletAddress, err := helpers.GetWalletAddress(issuer)
	if err != nil {
		return err
	}

	var cryptocurrency float64
	database.DB.Model(&model.CryptoWallet{}).Where("wallet_address = ? AND crypto_name = ?", WalletAddress, cryptoName).Pluck("amount", &cryptocurrency)
	if cryptocurrency < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid number",
		})
	}

	CryptoWallet(issuer, cryptoName, cryptoPrice, amountToSell, "sell")
	TransactionBalance(c, issuer, WalletAddress, totalProfit, "Sales", "Crypto Sales")
	TransactionCryptos(c, issuer, WalletAddress, cryptoPrice, cryptoName, amountToSell, "Sell")

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

	response := SellCryptoResponse{
		TotalProfit:      totalProfit,
		CryptoName:       cryptoName,
		AmountToSell:     amountToSell,
		Issuer:           issuer,
		UserBalance:      userBalance,
		UserBalanceAfter: totalBalance,
		Currency:         cryptocurrency,
	}

	return c.Status(200).JSON(response)
}

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

func TransactionListBalance(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuerFromContext(c)
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

	if err := database.DB.Create(&TransactionCryptos).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Transacton successful",
	})
}

func TransactionListCrypto(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuerFromContext(c)
	if err != nil {
		return err
	}

	var TransactionCryptos []model.TransactionCrypto
	if err := database.DB.Where("user_id = ?", issuer).Find(&TransactionCryptos).Error; err != nil {
		return err
	}

	return c.JSON(TransactionCryptos)
}

func CryptoWallet(User string, CryptoName string, CryptoPrice float64, Amount float64, ProcessType string) error {
	var existingCryptoWallet model.CryptoWallet
	WalletAddress, err := helpers.GetWalletAddress(User)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	UserInt, err := strconv.Atoi(User)

	result := database.DB.Preload("Wallet", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("wallet_address = ?", WalletAddress)
	}).Where("crypto_name = ?", CryptoName).First(&existingCryptoWallet)

	CryptoTotalPrice := CryptoPrice * Amount
	fmt.Println(CryptoName, CryptoPrice, Amount, ProcessType)

	newCryptoWallet := model.CryptoWallet{
		WalletID:         UserInt,
		CryptoName:       CryptoName,
		CryptoTotalPrice: CryptoTotalPrice,
		Amount:           Amount,
	}

	if ProcessType == "buy" {
		if result.Error == nil {
			existingCryptoWallet.Amount += Amount
			existingCryptoWallet.CryptoTotalPrice += CryptoTotalPrice
			database.DB.Save(&existingCryptoWallet)
		} else {
			database.DB.Create(&newCryptoWallet)
		}
	} else if ProcessType == "sell" {
		if result.Error == nil {
			existingCryptoWallet.Amount -= Amount
			existingCryptoWallet.CryptoTotalPrice -= CryptoTotalPrice
			database.DB.Save(&existingCryptoWallet)
		} else {
			database.DB.Create(&newCryptoWallet)
		}
	}

	return nil
}

func ListCryptoWallet(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuerFromContext(c)
	if err != nil {
		return err
	}

	var cryptoWallets []model.CryptoWallet

	database.DB.Preload("Wallet").Where("wallet_id = ?", issuer).Find(&cryptoWallets)

	type WalletListResponse struct {
		WalletAddress    string  `json:"wallet_address"`
		CryptoName       string  `json:"crypto_name"`
		Amount           float64 `json:"amount"`
		CryptoTotalPrice float64 `json:"crypto_total_price"`
	}

	WalletAddress, err := helpers.GetWalletAddress(issuer)
	if err != nil {
		return err
	}
	var response []WalletListResponse
	for _, cryptoWallet := range cryptoWallets {
		response = append(response, WalletListResponse{
			WalletAddress:    WalletAddress,
			CryptoName:       cryptoWallet.CryptoName,
			Amount:           cryptoWallet.Amount,
			CryptoTotalPrice: cryptoWallet.CryptoTotalPrice,
		})
	}
	return c.JSON(response)
}
