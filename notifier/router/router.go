package router

import (
	controller "Notifier/controllers"
	"Notifier/middleware"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	InitializeRouter(app)
}

func InitializeRouter(app *fiber.App) {
	app.Use(middleware.Logger())
	app.Post("/verfication", controller.Verification)
	app.Get("/verification-table", controller.PrintVerificationTable)

}
