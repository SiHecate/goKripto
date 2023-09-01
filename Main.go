package main

import (
	"gokripto/Database"
	"gokripto/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	Database.Connect()

	app := fiber.New()
	routes.Setup(app)
	app.Listen(":3000")

}
