package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Updatable
func timeoutHandler(timeoutDuration time.Duration) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ch := make(chan struct{})

		go func() {
			defer close(ch)
			err := c.Next()
			if err != nil {
				log.Println("İstek işlenirken hata oluştu:", err)
			}
		}()

		Loading := true

		for Loading {
			select {
			case <-time.After(timeoutDuration):
				log.Println(c.Route().Path)
				return c.Status(fiber.StatusRequestTimeout).JSON(fiber.Map{
					"message": "Timeout. Endpoint operation successful",
				})
			case <-ch:
				Loading = false
			}
		}
		return nil
	}
}
