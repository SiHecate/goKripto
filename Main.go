package main

import (
	"gokripto/Database"
	"gokripto/controllers"
	"gokripto/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	Database.Connect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)
	app.Static("/", "./static") // Bu satır, statik dosyaların sunulacağı klasörü belirtir
	app.Get("/get-exchange-rate", controllers.AddCryptoData)
	app.Get("/update-exchange-rate", controllers.UpdateCryptoData)

	app.Listen(":3000")
}
