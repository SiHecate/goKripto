package main

import (
	consume "Notifier/Consume"
	"Notifier/database"
	router "Notifier/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.Connect()
	app := fiber.New()

	router.Setup(app)

	go func() {
		consume.ConsumeMessages()
	}()

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n${body}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
	}))

	app.Listen(":8082")
}
