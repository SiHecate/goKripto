package controllers

import (
	"encoding/json"
	"fmt"
	model "gokripto/Model"
	"gokripto/database"
	helpers "gokripto/helpers"
	"net/http"
	"strconv"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ControllersApp interface {
	AddAllCryptoData(ws *websocket.Conn) error
	ListAllCryptos(c *fiber.Ctx) error
	AccountBalance(c *fiber.Ctx) error
	AddBalanceCrypto(c *fiber.Ctx) error
	BuyCryptos(c *fiber.Ctx) error
	SellCryptos(c *fiber.Ctx) error
	TransactionListBalance(c *fiber.Ctx) error
	TransactionListCrypto(c *fiber.Ctx) error
	ListCryptoWallet(c *fiber.Ctx) error
}

type WalletResponse struct {
	WalletAddress string  `json:"wallet_address"`
	Username      string  `json:"username"`
	Balance       float64 `json:"balance"`
}

type AddAllCryptoDataResponse struct {
	Message string `json:"message"`
}

type ListAllCryptoResponse struct {
	CryptoList []ListAllCrypto `json:"crypto_list"`
}

type ListAllCrypto struct {
	CryptoSymbol string  `json:"crypto_symbol"`
	CryptoName   string  `json:"crypto_name"`
	CryptoPrice  float64 `json:"crypto_price"`
}

type AccountBalanceResponse struct {
	WalletAddress string  `json:"wallet_address"`
	Username      string  `json:"username"`
	Balance       float64 `json:"balance"`
}

type AddBalanceResponse struct {
	Issuer           string  `json:"issuer"`
	AvailableBalance float64 `json:"available_balance"`
	TotalBalance     float64 `json:"total_balance"`
}

type BuyCryptoResponse struct {
	TotalCost     float64 `json:"total_cost"`
	CryptoName    string  `json:"crypto_name"`
	AmountToBuy   float64 `json:"amount_to_buy"`
	Issuer        string  `json:"issuer"`
	UserBalance   float64 `json:"user_balance"`
	UserBalanceAB float64 `json:"user_balance_after_buy"`
}

type SellCryptoResponse struct {
	TotalProfit      float64 `json:"total_profit"`
	CryptoName       string  `json:"crypto_name"`
	AmountToSell     float64 `json:"amount_to_sell"`
	Issuer           string  `json:"issuer"`
	UserBalance      float64 `json:"user_balance"`
	UserBalanceAfter float64 `json:"user_balance_after_sell"`
}

type TransactionResponse struct {
	UserID        string  `json:"user_id"`
	WalletAddres  string  `json:"wallet_address"`
	BalanceAmount float64 `json:"price"`
	Type          string  `json:"type"`
	TypeInfo      string  `json:"type_info"`
}

type TransactionResponseCrypto struct {
	UserID       string  `json:"user_id"`
	WalletAddres string  `json:"wallet_address"`
	Price        float64 `json:"price"`
	CryptoName   string  `json:"crypto_name"`
	Amount       float64 `json:"amount"`
	Type         string  `json:"type"`
}

type WalletListResponse struct {
	WalletAddress    string  `json:"wallet_address"`
	CryptoName       string  `json:"crypto_name"`
	Amount           float64 `json:"amount"`
	CryptoTotalPrice float64 `json:"crypto_total_price"`
}

type Status400Response struct {
	Message string `json:"StatusBadRequest"`
}

type Status401Reponse struct {
	Message string `json:"StatusUnauthorized"`
}

type Status404Response struct {
	Message string `json:"StatusNotFound"`
}

// AddAllCryptoData
// @Summary Fetch and Add Crypto Data using WebSocket
// @Description Fetch cryptocurrency data from an external API and add/update it in the database using WebSocket.
// @Tags Cryptocurrencies
// @Security ApiKeyAuth
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} Status400Response
// @Router /websocket [get]
func AddAllCryptoData(ws *websocket.Conn) error {
	apiURL := "https://api.coincap.io/v2/assets"

	response, err := http.Get(apiURL)
	if err != nil {
		errorMessage := fiber.Map{
			"message": "Failed to fetch crypto data from the API.",
			"error":   err.Error(),
		}
		helpers.SendJSONError(ws, errorMessage)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		errorMessage := fiber.Map{
			"message": "HTTP Error: " + strconv.Itoa(response.StatusCode),
		}
		helpers.SendJSONError(ws, errorMessage)
		return nil
	}

	var cryptoAPIResponse model.CryptoAPIResponse
	err = json.NewDecoder(response.Body).Decode(&cryptoAPIResponse)
	if err != nil {
		errorMessage := fiber.Map{
			"message": "Failed to parse API response.",
			"error":   err.Error(),
		}
		helpers.SendJSONError(ws, errorMessage)
		return err
	}

	for _, crypto := range cryptoAPIResponse.Data {
		if err := createOrUpdateCrypto(crypto.Symbol, crypto.Name, crypto.PriceUsd); err != nil {
			errorMessage := fiber.Map{
				"message": "Failed to add or update crypto data in the database.",
				"error":   err.Error(),
			}
			helpers.SendJSONError(ws, errorMessage)
			return err
		}
	}

	successMessage := fiber.Map{
		"message": "All crypto data added or updated successfully",
	}
	helpers.SendJSONResponse(ws, successMessage)

	return nil
}

func createOrUpdateCrypto(symbol string, name string, price string) error {
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return err
	}

	var crypto model.Crypto
	result := database.Conn.Where("symbol = ?", symbol).First(&crypto)
	switch {
	case result.Error != nil:
		crypto := model.Crypto{
			Symbol: symbol,
			Name:   name,
			Price:  priceFloat,
		}
		if err := database.Conn.Create(&crypto).Error; err != nil {
			return err
		}
	case result.Error != nil:
		return err

	default:
		crypto.Price = priceFloat
		database.Conn.Save(&crypto)
	}

	return nil
}

// ListAllCryptos
// @Summary List all cryptocurrencies
// @Description Get a list of all cryptocurrencies with their symbols, names, and prices.
// @Tags Cryptocurrencies
// @Accept json
// @Produce json
// @Success 200 {object} []ListAllCryptoResponse
// @Router /crypto/cryptoList [get]
func ListAllCryptos(c *fiber.Ctx) error {
	cryptos, err := model.GetAllCryptos(database.Conn)
	if err != nil {
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

// AccountBalance
// @Summary Get User Account Balance
// @Description Get the account balance of the authenticated user.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} WalletResponse
// @Failure 400 {object} Status400Response
// @Failure 401 {object} Status401Reponse
// @Failure 404 {object} Status404Response
// @Router /user/balance [get]
func AccountBalance(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuer(c)
	if err != nil {
		return err
	}

	WalletAddress, err := model.GetWalletAddressByIssuer(database.Conn, issuer)
	if err != nil {
		return err
	}

	wallet, err := model.GetWalletbyWalletAddress(database.Conn, WalletAddress)
	if err != nil {
		return err
	}

	user, err := model.GetUserByIssuer(database.Conn, issuer)
	if err != nil {
		return err
	}

	walletAddress, err := model.GetWalletAddress(database.Conn, issuer)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Wallet not found",
		})
	}

	type WalletResponse struct {
		WalletAddress string  `json:"wallet_address"`
		Username      string  `json:"username"`
		Balance       float64 `json:"balance"`
	}

	response := WalletResponse{
		WalletAddress: walletAddress,
		Username:      user.Name,
		Balance:       wallet.Balance,
	}

	return c.JSON(response)
}

// AddBalanceCrypto
// @Summary Add Balance to Crypto Wallet
// @Description Add a specified amount to the user's cryptocurrency wallet balance.
// @Tags User
// @Accept json
// @Produce json
// @Param add_balance body float64 true "Amount to add to the balance"
// @Security ApiKeyAuth
// @Success 200 {object} AddBalanceResponse
// @Failure 400 {object} Status400Response
// @Router /user/addBalance [post]
func AddBalanceCrypto(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuer(c)
	if err != nil {
		return err
	}

	var data struct {
		AddBalance float64 `json:"add_balance"`
	}

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	addBalance := data.AddBalance

	walletAddress, err := model.GetWalletAddressByIssuer(database.Conn, issuer)
	if err != nil {
		return err
	}

	var availableBalance float64
	if err := database.Conn.Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &availableBalance).Error; err != nil {
		return err
	}

	totalBalance := addBalance + availableBalance

	if err := database.Conn.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	TransactionBalance(c, issuer, walletAddress, addBalance, "Deposit", "Balance Adding")

	type AddBalanceResponse struct {
		Issuer           string  `json:"issuer"`
		AvailableBalance float64 `json:"available_balance"`
		TotalBalance     float64 `json:"added_balance"`
	}

	response := AddBalanceResponse{
		Issuer:           issuer,
		AvailableBalance: availableBalance,
		TotalBalance:     totalBalance,
	}

	return c.JSON(response)
}

// BuyCryptos
// @Summary Buy Cryptocurrency
// @Description Buy a specified amount of cryptocurrency.
// @Tags Cryptocurrencies
// @Accept json
// @Produce json
// @Param cryptoName body string true "Name of the crypto currency to buy"
// @Param amountToBuy body float64 true "Amount of crypto currency to buy"
// @Security ApiKeyAuth
// @Success 200 {object} BuyCryptoResponse
// @Failure 400 {object} Status400Response
// @Router /crypto/cryptoBuy [post]
func BuyCryptos(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuer(c)
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
	if err := database.Conn.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice).Error; err != nil {
		return err
	}

	var userBalance float64
	if err := database.Conn.Model(&model.Wallet{}).Where("user_id = ?", issuer).Pluck("balance", &userBalance).Error; err != nil {
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

	if err := database.Conn.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	var cryptoID uint
	if err := database.Conn.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID).Error; err != nil {
		return err
	}

	modelWallet := model.Wallet{}
	database.Conn.Where("user_id = ?", issuer).Find(&modelWallet)

	WalletAddress, err := model.GetWalletAddressByIssuer(database.Conn, issuer)
	if err != nil {
		return err
	}

	CryptoWallet(issuer, cryptoName, cryptoPrice, amountToBuy, "buy")
	TransactionBalance(c, issuer, WalletAddress, totalCost, "Purchase", "Crypto Purchase")
	TransactionCryptos(c, issuer, WalletAddress, cryptoPrice, cryptoName, amountToBuy, "Buy")

	type BuyCryptoResponse struct {
		TotalCost     float64 `json:"total_cost"`
		CryptoName    string  `json:"crypto_name"`
		AmountToBuy   float64 `json:"amount_to_buy"`
		Issuer        string  `json:"issuer"`
		UserBalance   float64 `json:"user_balance"`
		UserBalanceAB float64 `json:"user_balance_after_buy"`
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

// SellCryptos
// @Summary Sell Cryptocurrency
// @Description Sell a specified amount of cryptocurrency.
// @Tags Cryptocurrencies
// @Accept json
// @Produce json
// @Param cryptoName body string true "Name of the crypto currency to sell"
// @Param amountToSell body float64 true "Amount of crypto currency to sell"
// @Security ApiKeyAuth
// @Success 200 {object} SellCryptoResponse
// @Failure 400 {object} Status400Response
// @Router /crypto/cryptoSell [post]
func SellCryptos(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuer(c)
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
	if err := database.Conn.Model(&model.Wallet{}).Where("id = ?", issuer).Pluck("balance", &userBalance).Error; err != nil {
		return err
	}

	var cryptoPrice float64
	database.Conn.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("price", &cryptoPrice)

	totalProfit := cryptoPrice * amountToSell
	totalBalance := userBalance + totalProfit
	if err := database.Conn.Model(&model.Wallet{}).Where("user_id = ?", issuer).Update("balance", totalBalance).Error; err != nil {
		return err
	}

	var cryptoID uint
	if err := database.Conn.Model(&model.Crypto{}).Where("name = ?", cryptoName).Pluck("id", &cryptoID).Error; err != nil {
		return err
	}

	modelWallet := model.Wallet{}
	database.Conn.Where("user_id = ?", issuer).Find(&modelWallet)

	WalletAddress, err := model.GetWalletAddressByIssuer(database.Conn, issuer)
	if err != nil {
		return err
	}

	CryptoWallet(issuer, cryptoName, cryptoPrice, amountToSell, "sell")
	TransactionBalance(c, issuer, WalletAddress, totalProfit, "Sales", "Crypto Sales")
	TransactionCryptos(c, issuer, WalletAddress, cryptoPrice, cryptoName, amountToSell, "Sell")

	type SellCryptoResponse struct {
		TotalProfit      float64 `json:"total_profit"`
		CryptoName       string  `json:"crypto_name"`
		AmountToSell     float64 `json:"amount_to_sell"`
		Issuer           string  `json:"issuer"`
		UserBalance      float64 `json:"user_balance"`
		UserBalanceAfter float64 `json:"user_balance_after_sell"`
	}

	response := SellCryptoResponse{
		TotalProfit:      totalProfit,
		CryptoName:       cryptoName,
		AmountToSell:     amountToSell,
		Issuer:           issuer,
		UserBalance:      userBalance,
		UserBalanceAfter: totalBalance,
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
	}

	if err := database.Conn.Create(&TransactionBalance).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Transaction successful",
	})
}

// TransactionListBalance
// @Summary List Balance Transactions
// @Description Get a list of balance transactions
// @Tags Transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} TransactionResponse
// @Failure 400 {object} Status400Response
// @Router /transaction/balanceTransaction [get]
func TransactionListBalance(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuer(c)
	if err != nil {
		return err
	}

	var TransactionBalance []model.TransactionBalance
	if err := database.Conn.Where("user_id = ?", issuer).Find(&TransactionBalance).Error; err != nil {
		return err
	}

	type TransactionResponse struct {
		UserID        string  `json:"user_id"`
		WalletAddres  string  `json:"wallet_address"`
		BalanceAmount float64 `json:"price"`
		Type          string  `json:"type"`
		TypeInfo      string  `json:"type_info"`
	}

	var response []TransactionResponse

	for _, transaction := range TransactionBalance {
		response = append(response, TransactionResponse{
			UserID:        issuer,
			WalletAddres:  transaction.WalletAddress,
			BalanceAmount: transaction.BalanceAmount,
			Type:          transaction.Type,
			TypeInfo:      transaction.TypeInfo,
		})
	}

	return c.JSON(response)
}

func TransactionCryptos(c *fiber.Ctx, UserID string, WalletAddres string, price float64, cryptoname string, amount float64, transactionType string) error {
	TransactionCryptos := model.TransactionCrypto{
		UserID:       UserID,
		WalletAddres: WalletAddres,
		Price:        price,
		CryptoName:   cryptoname,
		Amount:       amount,
		Type:         transactionType,
	}

	if err := database.Conn.Create(&TransactionCryptos).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Transacton successful",
	})
}

// TransactionListCrypto
// @Summary List Crypto Transactions
// @Description Get a list of cryptocurrency transactions
// @Tags Transactions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} TransactionResponse
// @Failure 400 {object} Status400Response
// @Router /transaction/cryptoTransaction [get]
func TransactionListCrypto(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuer(c)
	if err != nil {
		return err
	}

	var TransactionCrypto []model.TransactionCrypto
	if err := database.Conn.Where("user_id = ?", issuer).Find(&TransactionCrypto).Error; err != nil {
		return err
	}

	type TransactionResponse struct {
		UserID       string  `json:"user_id"`
		WalletAddres string  `json:"wallet_address"`
		Price        float64 `json:"price"`
		CryptoName   string  `json:"crypto_name"`
		Amount       float64 `json:"amount"`
		Type         string  `json:"type"`
	}

	var response []TransactionResponse

	for _, transaction := range TransactionCrypto {
		response = append(response, TransactionResponse{
			UserID:       issuer,
			WalletAddres: transaction.WalletAddres,
			Price:        transaction.Price,
			CryptoName:   transaction.CryptoName,
			Amount:       transaction.Amount,
			Type:         transaction.Type,
		})
	}

	return c.JSON(response)
}

func CryptoWallet(User string, CryptoName string, CryptoPrice float64, Amount float64, ProcessType string) error {
	var existingCryptoWallet model.CryptoWallet
	WalletAddress, err := model.GetWalletAddressByIssuer(database.Conn, User)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	UserInt, err := strconv.Atoi(User)

	result := database.Conn.Preload("Wallet", func(tx *gorm.DB) *gorm.DB {
		return tx.Where("wallet_address = ?", WalletAddress)
	}).Where("crypto_name = ?", CryptoName).First(&existingCryptoWallet)

	CryptoTotalPrice := CryptoPrice * Amount

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
			database.Conn.Save(&existingCryptoWallet)
		} else {
			database.Conn.Create(&newCryptoWallet)
		}
	} else if ProcessType == "sell" {
		if result.Error == nil {
			existingCryptoWallet.Amount -= Amount
			existingCryptoWallet.CryptoTotalPrice -= CryptoTotalPrice
			database.Conn.Save(&existingCryptoWallet)
		} else {
			database.Conn.Create(&newCryptoWallet)
		}
	}

	return nil
}

// ListCryptoWallet
// @Summary List Crypto Wallet
// @Description Get a list of cryptocurrencies in a user's wallet.
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} WalletListResponse
// @Failure 400 {object} Status400Response
// @Router /user/cryptoWallet [get]
func ListCryptoWallet(c *fiber.Ctx) error {
	issuer, err := helpers.GetIssuer(c)
	if err != nil {
		return err
	}

	WalletAddress, err := model.GetWalletAddressByIssuer(database.Conn, issuer)
	if err != nil {
		return err
	}

	WalletID, err := model.GetCryptoWallet(database.Conn, issuer)
	if err != nil {
		return err
	}

	type WalletListResponse struct {
		WalletAddress    string  `json:"wallet_address"`
		CryptoName       string  `json:"crypto_name"`
		Amount           float64 `json:"amount"`
		CryptoTotalPrice float64 `json:"crypto_total_price"`
	}

	var response []WalletListResponse
	for _, cryptoWallet := range WalletID {
		response = append(response, WalletListResponse{
			WalletAddress:    WalletAddress,
			CryptoName:       cryptoWallet.CryptoName,
			Amount:           cryptoWallet.Amount,
			CryptoTotalPrice: cryptoWallet.CryptoTotalPrice,
		})
	}
	return c.JSON(response)
}
