package routes

import (
	AuthController "gokripto/controllers/AuthController"
	CryptoControllers "gokripto/controllers/CryptoController"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	//Post

	app.Post("/api/register", AuthController.Register)         // test done
	app.Post("/api/login", AuthController.Login)               // test done
	app.Post("/api/logout", AuthController.Logout)             // test done
	app.Post("/api/cryptoBuy", CryptoControllers.BuyCryptos)   // test done
	app.Post("/api/cryptoSell", CryptoControllers.SellCryptos) // test done
	app.Post("/api/addBalance", CryptoControllers.AddBalanceCrypto)

	//Get

	app.Get("/api/user", AuthController.User)
	app.Get("/api/balance", CryptoControllers.AccountBalance)
	app.Get("/api/cryptoAdd", CryptoControllers.AddCryptoData)
	app.Get("/api/cryptoUpdate", CryptoControllers.UpdateCryptoData)
	app.Get("/api/cryptoList", CryptoControllers.ListAllCryptos)
	app.Get("/api/listcryptowallet", CryptoControllers.ListCryptoWallet)

}
