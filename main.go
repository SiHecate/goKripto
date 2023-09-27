package main

import (
	"gokripto/database"
	websocket "gokripto/pkg/controllers"
	routes "gokripto/routes"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()

	//database.DownTables()

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
