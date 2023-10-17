package main

import (
	consume "Notifier/Consume"
	"Notifier/database"
	"Notifier/logger"

	router "Notifier/router"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	database.Connect()

	logger, _ := logger.InitLogger()

	app := fiber.New()

	go func() {
		consume.ConsumeMessages()
	}()

	app.Use(func(c *fiber.Ctx) error {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()

		logger.Info("HTTP log",
			zap.String("Method", c.Method()),
			zap.String("Path", c.Path()),
			zap.Int("Status", c.Response().StatusCode()),
			zap.Duration("Latency", endTime.Sub(startTime)),
		)

		return nil
	})

	router.Setup(app)

	app.Listen(":8082")
}
