package routes

import (
	AuthController "gokripto/controllers/auth"
	CryptoControllers "gokripto/controllers/crypto"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	InitializeRouter(app)
}

func InitializeRouter(app *fiber.App) {
	user := app.Group("/user")
	user.Post("/register", AuthController.Register)
	user.Post("/login", AuthController.Login)
	user.Get("/user", AuthController.User)
	user.Post("/logout", AuthController.Logout)
	user.Post("/cryptoWallet", CryptoControllers.ListCryptoWallet)
	user.Get("/balance", CryptoControllers.AccountBalance)
	user.Post("/addBalance", CryptoControllers.AddbalanceCrypto)

	crypto := app.Group("/crypto")
	crypto.Post("/cryptoBuy", CryptoControllers.BuyCryptos)
	crypto.Post("/cryptoSell", CryptoControllers.SellCryptos)
	crypto.Get("/cryptoList", CryptoControllers.ListAllCryptos)

	transaction := app.Group("/transaction")
	transaction.Get("/cryptoTransaction", CryptoControllers.TransactionListCrypto)
	transaction.Get("/balanceTransaction", CryptoControllers.TransactionListBalance)
}
