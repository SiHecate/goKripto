package router

import (
	AuthController "cryptoApp/controllers/auth"
	CryptoControllers "cryptoApp/controllers/crypto"
	"cryptoApp/middlewares"
	"cryptoApp/produce"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/gofiber/swagger"
)

func Setup(app *fiber.App) {
	InitializeRouter(app)
}

const SecretKey = "secret"

func InitializeRouter(app *fiber.App) {

	// Hello World
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	// Prometheus
	prometheus := fiberprometheus.New("my-service-name")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// Verification endpoint
	verification := app.Group("/ver")
	verification.Use(middlewares.JWTMiddleware())
	verification.Get("/verification_code", produce.VerficationCode)

	// Swagger
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	// Authentication
	auth := app.Group("/auth")
	auth.Use(middlewares.Logger())
	auth.Post("/register", AuthController.Register)
	auth.Post("/login", AuthController.Login)

	// User
	user := app.Group("/user")
	user.Use(middlewares.JWTMiddleware())
	user.Use(middlewares.Validiation())
	user.Use(middlewares.Logger())
	user.Post("/logout", AuthController.Logout)
	user.Post("/addBalance", CryptoControllers.AddBalanceCrypto)
	user.Get("/user", AuthController.User)
	user.Get("/cryptoWallet", CryptoControllers.ListCryptoWallet)
	user.Get("/balance", CryptoControllers.AccountBalance)

	// Crypto
	crypto := app.Group("/crypto")
	crypto.Use(middlewares.JWTMiddleware())
	crypto.Use(middlewares.Validiation())
	crypto.Use(middlewares.Logger())
	crypto.Post("/cryptoBuy", CryptoControllers.BuyCryptos)
	crypto.Post("/cryptoSell", CryptoControllers.SellCryptos)
	crypto.Get("/cryptoList", CryptoControllers.ListAllCryptos)

	// Transaction
	transaction := app.Group("/transaction")
	transaction.Use(middlewares.JWTMiddleware())
	transaction.Use(middlewares.Validiation())
	transaction.Use(middlewares.Logger())
	transaction.Get("/cryptoTransaction", CryptoControllers.TransactionListCrypto)
	transaction.Get("/balanceTransaction", CryptoControllers.TransactionListBalance)

}
