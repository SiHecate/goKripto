package routes

import (
	AuthController "gokripto/controllers/auth"
	CryptoControllers "gokripto/controllers/crypto"
	"gokripto/middlewares"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

func Setup(app *fiber.App) {
	InitializeRouter(app)
}

const SecretKey = "secret"

func InitializeRouter(app *fiber.App) {
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	auth := app.Group("/auth")
	auth.Post("/register", AuthController.Register)
	auth.Post("/login", AuthController.Login)

	user := app.Group("/user")
	user.Use(middlewares.JWTMiddleware())
	user.Post("/logout", AuthController.Logout)
	user.Post("/addBalance", CryptoControllers.AddBalanceCrypto)
	user.Get("/user", AuthController.User)
	user.Get("/cryptoWallet", CryptoControllers.ListCryptoWallet)
	user.Get("/balance", CryptoControllers.AccountBalance)

	crypto := app.Group("/crypto")
	crypto.Use(middlewares.JWTMiddleware())
	crypto.Post("/cryptoBuy", CryptoControllers.BuyCryptos)
	crypto.Post("/cryptoSell", CryptoControllers.SellCryptos)
	crypto.Get("/cryptoList", CryptoControllers.ListAllCryptos)

	transaction := app.Group("/transaction")
	transaction.Use(middlewares.JWTMiddleware())
	transaction.Get("/cryptoTransaction", CryptoControllers.TransactionListCrypto)
	transaction.Get("/balanceTransaction", CryptoControllers.TransactionListBalance)

}
