package router

import (
	controller "Notifier/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	InitializeRouter(app)
}

func InitializeRouter(app *fiber.App) {
	app.Post("/verfication", controller.Verification)
	app.Get("/verification-table", controller.PrintVerificationTable)

}
