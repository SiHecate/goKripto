package main

import (
	websocket "gokripto/controllers"
	"gokripto/database"
	routes "gokripto/routes"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/gofiber/swagger"

	_ "gokripto/docs"
)

// @title           Go Crypto
// @version         1.0
// @description     Crypto currency app.

// @contact.name   API Support
// @contact.url    https://github.com/SiHecate

// @host      localhost:3000

// @securityDefinitions.basic  BasicAuth
func main() {
	database.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

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
