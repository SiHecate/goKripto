package middlewares

import (
	"cryptoApp/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func Logger() fiber.Handler {
	logger, _ := logger.InitLogger()
	defer logger.Sync()

	return func(c *fiber.Ctx) error {
		startTime := time.Now()
		if err := c.Next(); err != nil {
			return err
		}
		endTime := time.Now()

		logger.Info("HTTP Log",
			zap.String("Method", c.Method()),
			zap.String("Path", c.Path()),
			zap.Int("Status", c.Response().StatusCode()),
			zap.Duration("Latency", endTime.Sub(startTime)),
		)

		return nil
	}
}
