package websocket

import (
	CryptoController "gokripto/controllers/CryptoController"
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

	app.Get("/ws/:id", websocket.New(func(ws *websocket.Conn) {
		log.Printf("WebSocket port open on: %s", ws.Params("id"))
		// Yenileme hızı
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// Websocket
		for range ticker.C {
			CryptoController.UpdateWSCryptoData(ws)
		}

		log.Printf("WebSocket port off!: %s", ws.Params("id"))
	}))

	app.Listen(":3000")
}
