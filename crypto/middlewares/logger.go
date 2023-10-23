package middlewares

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	zaploki "github.com/paul-milne/zap-loki"
	"go.uber.org/zap"
)

// Loki sender
func InitLogger() (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()
	loki := zaploki.New(context.Background(), zaploki.Config{
		Url:          "http://loki:3100",
		BatchMaxSize: 1000,
		BatchMaxWait: 10 * time.Second,
		Labels:       map[string]string{"app": "cryptoApp"},
	})

	return loki.WithCreateLogger(zapConfig)
}

func Logger() fiber.Handler {
	logger, _ := InitLogger()
	defer logger.Sync()

	return func(c *fiber.Ctx) error {

		// latency calculating
		startTime := time.Now()
		if err := c.Next(); err != nil {
			return err
		}
		endTime := time.Now()
		// latency calculating

		logger.Info("HTTP Log",
			zap.String("Method", c.Method()),
			zap.String("Path", c.Path()),
			zap.Int("Status", c.Response().StatusCode()),
			zap.Duration("Latency", endTime.Sub(startTime)),
		)

		return nil
	}
}
