package main

import (
	"gokripto/Database"
	websocket "gokripto/Websocket"
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
	go websocket.StartWebSocket(app)
	app.Listen(":3000")
}
