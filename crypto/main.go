package main

import (
	"cryptoApp/database"
	router "cryptoApp/router"
	websocket "cryptoApp/router"
	"sync"

	_ "cryptoApp/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title           Go Cryptos
// @version         1.0
// @description     Crypto currency app.
// @contact.name   API Support
// @contact.url    https://github.com/SiHecate
// @host      localhost:8080
// @securityDefinitions.basic  BasicAuth
func main() {
	database.Connect()
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format:     "${time} ${status} - ${method} ${path}\n${body}\n",
		TimeFormat: "02-Jan-2006 15:04:05",
	}))

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	router.Setup(app)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		websocket.StartWebSocket(app)
		wg.Done()
	}()
	wg.Wait()

	app.Listen(":8080")
}
