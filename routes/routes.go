package routes

import (
	AuthController "gokripto/controllers/AuthController"
	CryptoControllers "gokripto/controllers/CryptoController"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	timeoutDuration := 2000 * time.Millisecond

	timeoutHandler := func(c *fiber.Ctx) error {
		ch := make(chan struct{})

		go func() {
			defer close(ch)
			err := c.Next()
			if err != nil {
				log.Println("İstek işlenirken hata oluştu:", err)
			}
		}()

		Loading := true

		for Loading {
			select {
			case <-time.After(timeoutDuration):
				log.Println(c.Route().Path)
				return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
					"message": "Timeout. Endpoint operation successful",
				})
			case <-ch:
				Loading = false
			}
		}
		return nil
	}

	// Post
	app.Post("/api/register", timeoutHandler, AuthController.Register)              // Done
	app.Post("/api/login", timeoutHandler, AuthController.Login)                    // Done
	app.Post("/api/logout", timeoutHandler, AuthController.Logout)                  // Done
	app.Post("/api/cryptoBuy", timeoutHandler, CryptoControllers.BuyCryptos)        // Done
	app.Post("/api/cryptoSell", timeoutHandler, CryptoControllers.SellCryptos)      // Done
	app.Post("/api/addBalance", timeoutHandler, CryptoControllers.AddBalanceCrypto) // Done

	// Get
	app.Get("/api/CryptoTransactionHistory", timeoutHandler, CryptoControllers.TransactionListCrypto)   // Done
	app.Get("/api/BalanceTransactionHistory", timeoutHandler, CryptoControllers.TransactionListBalance) // Done
	app.Get("/api/user", timeoutHandler, AuthController.User)                                           // Done
	app.Get("/api/balance", timeoutHandler, CryptoControllers.AccountBalance)                           // Done
	app.Get("/api/cryptoAdd", CryptoControllers.AddAllCryptoData)                                       // ?
	app.Get("/api/cryptoList", timeoutHandler, CryptoControllers.ListAllCryptos)                        // Done
	app.Get("/api/listcryptowallet", timeoutHandler, CryptoControllers.ListCryptoWallet)                // Done
}
