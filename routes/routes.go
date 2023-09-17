package routes

import (
	AuthController "gokripto/controllers/AuthController"
	CryptoControllers "gokripto/controllers/CryptoController"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	timeoutDuration := 650 * time.Millisecond

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
	app.Post("/api/register", timeoutHandler, AuthController.Register)
	app.Post("/api/login", timeoutHandler, AuthController.Login)
	app.Post("/api/logout", timeoutHandler, AuthController.Logout)
	app.Post("/api/cryptoBuy", timeoutHandler, CryptoControllers.BuyCryptos)
	app.Post("/api/cryptoSell", timeoutHandler, CryptoControllers.SellCryptos)
	app.Post("/api/addBalance", timeoutHandler, CryptoControllers.AddBalanceCrypto)

	// Get
	app.Get("/api/user", timeoutHandler, AuthController.User)
	app.Get("/api/balance", timeoutHandler, CryptoControllers.AccountBalance)
	app.Get("/api/cryptoAdd", CryptoControllers.AddCryptoData)
	app.Get("/api/cryptoUpdate", timeoutHandler, CryptoControllers.UpdateCryptoData)
	app.Get("/api/cryptoList", timeoutHandler, CryptoControllers.ListAllCryptos)
	app.Get("/api/listcryptowallet", timeoutHandler, CryptoControllers.ListCryptoWallet)
}
