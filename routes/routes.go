package routes

import (
	AuthController "gokripto/pkg/controllers/auth"
	CryptoControllers "gokripto/pkg/controllers/crypto"
	Middleware "gokripto/pkg/middlewares"

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	timeoutDuration := 2000 * time.Millisecond

	setupAllRoutes(app, timeoutHandler(timeoutDuration))
}

func timeoutHandler(timeoutDuration time.Duration) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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
}

func setupAllRoutes(app *fiber.App, timeoutHandler func(*fiber.Ctx) error) {
	setupPostRoutes(app, timeoutHandler)
	setupGetRoutes(app, timeoutHandler)

}

func setupPostRoutes(app *fiber.App, timeoutHandler func(*fiber.Ctx) error) {
	app.Post("/api/register", timeoutHandler, AuthController.Register)
	app.Post("/api/login", timeoutHandler, AuthController.Login)
	app.Post("/api/logout", timeoutHandler, AuthController.Logout)
	app.Post("/api/cryptoBuy", Middleware.GetIssuer, timeoutHandler, CryptoControllers.BuyCryptos)
	app.Post("/api/cryptoSell", Middleware.GetIssuer, timeoutHandler, CryptoControllers.SellCryptos)
	app.Post("/api/addBalance", Middleware.GetIssuer, timeoutHandler, CryptoControllers.AddBalanceCrypto)
}

func setupGetRoutes(app *fiber.App, timeoutHandler func(*fiber.Ctx) error) {
	app.Get("/api/CryptoTransactionHistory", Middleware.GetIssuer, timeoutHandler, CryptoControllers.TransactionListCrypto)
	app.Get("/api/BalanceTransactionHistory", Middleware.GetIssuer, timeoutHandler, CryptoControllers.TransactionListBalance)
	app.Get("/api/user", timeoutHandler, AuthController.User)
	app.Get("/api/balance", Middleware.GetIssuer, timeoutHandler, CryptoControllers.AccountBalance)
	app.Get("/api/cryptoList", timeoutHandler, CryptoControllers.ListAllCryptos)
	app.Get("/api/listcryptowallet", Middleware.GetIssuer, timeoutHandler, CryptoControllers.ListCryptoWallet)
}
