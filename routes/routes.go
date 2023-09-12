package routes

import (
	AuthController "gokripto/controllers/AuthController"
	CryptoControllers "gokripto/controllers/CryptoController"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/api/register", AuthController.Register)
	app.Post("/api/login", AuthController.Login)
	app.Get("/api/user", AuthController.User)
	app.Get("/api/balance", CryptoControllers.AccountBalance)
	app.Post("/api/logout", AuthController.Logout)
	app.Post("/api/cryptoBuy", CryptoControllers.BuyCryptos)
	app.Post("/api/cryptoSell", CryptoControllers.SellCryptos)
	app.Get("/api/cryptoAdd", CryptoControllers.AddCryptoData)
	app.Get("/api/cryptoUpdate", CryptoControllers.UpdateCryptoData)
	app.Get("/api/cryptoList", CryptoControllers.ListAllCryptos)
	app.Get("/api/listcryptowallet", CryptoControllers.ListCryptoWallet)
	app.Post("/api/addBalance", CryptoControllers.AddBalanceCrypto)
}
