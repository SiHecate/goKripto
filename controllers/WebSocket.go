package controllers

import (
	controllers "gokripto/controllers/crypto"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func StartWebSocket(app *fiber.App) {
	app.Use("/websocket", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/websocket", websocket.New(func(ws *websocket.Conn) {
		log.Println("WebSocket port open")
		// Yenileme hızı
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		// Websocket
		for range ticker.C {
			controllers.AddAllCryptoData(ws.Conn)
		}

		log.Println("WebSocket port off")
	}))

	app.Listen(":3000")
}
