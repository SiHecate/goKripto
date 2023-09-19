package main

import (
	"gokripto/Database"
	websocket "gokripto/Websocket"
	CryptoControllers "gokripto/controllers/CryptoController"
	"gokripto/routes"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	CryptoControllers.CryptoBill()
	Database.Connect()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	routes.Setup(app)
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		websocket.StartWebSocket(app)
		wg.Done()
	}()
	wg.Wait()
	app.Listen(":3000")
}
