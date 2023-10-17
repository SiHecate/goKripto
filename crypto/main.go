package main

import (
	"cryptoApp/database"
	logger "cryptoApp/logger"
	router "cryptoApp/router"
	websocket "cryptoApp/router"
	"sync"
	"time"

	_ "cryptoApp/docs"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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

	logger, _ := logger.InitLogger()
	defer logger.Sync()

	app.Use(func(c *fiber.Ctx) error {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		logger.Info("HTTP Log",
			zap.String("Method", c.Method()),
			zap.String("Path", c.Path()),
			zap.Int("Status", c.Response().StatusCode()),
			zap.Duration("Latency", endTime.Sub(startTime)),
		)

		return nil
	})

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
