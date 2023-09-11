package websocket

import (
	"gokripto/Database"
	model "gokripto/Model"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func StartWebSocket(app *fiber.App) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:custom_id", websocket.New(func(c *websocket.Conn) {
		log.Printf("WebSocket port open!: %s", c.Params("custom_id"))
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()""

		var cryptos []model.Crypto
		Database.GetDB().Find(&cryptos)
		for range ticker.C {
			message := cryptos
			c.WriteJSON(message)
		}

		log.Printf("WebSocket port off!: %s", c.Params("id"))
	}))

	app.Listen(":3000")
}
