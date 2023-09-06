package routes

import (
	"gokripto/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
<<<<<<< HEAD
	app.Post("/api/cryptoBuy", controllers.BuyCryptos)
	app.Post("/api/cryptoSell", controllers.SellCryptos)
	app.Get("/api/cryptoAdd", controllers.AddCryptoData)
	app.Get("/api/cryptoUpdate", controllers.UpdateCryptoData)
	app.Get("/api/cryptoList", controllers.ListAllCryptos)
=======
	app.Post("/api/Buy", controllers.BuyCryptos)

>>>>>>> origin/main
}
