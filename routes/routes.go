package routes

import (
	AuthController "gokripto/controllers/auth"
	CryptoControllers "gokripto/controllers/crypto"
	"gokripto/middlewares"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	InitializeRouter(app)
	InitializePrometheus(app)
}

const SecretKey = "secret"

func InitializeRouter(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Use("/register", AuthController.Register)
	auth.Use("/login", AuthController.Login)

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
	transaction.Use(middlewares.JWTMiddleware()) // Add other middleware as needed
	transaction.Get("/cryptoTransaction", CryptoControllers.TransactionListCrypto)
	transaction.Get("/balanceTransaction", CryptoControllers.TransactionListBalance)
}

func InitializePrometheus(app *fiber.App) {
	prometheus := fiberprometheus.New("my-service-name")
	prometheus.RegisterAt(app, "/metrics")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})
	app.Post("/some", func(c *fiber.Ctx) error {
		return c.SendString("Welcome!")
	})
}
