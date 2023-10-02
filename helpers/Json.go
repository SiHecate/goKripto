package helpers

import (
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
)

func SendJSONError(ws *websocket.Conn, errorMessage fiber.Map) {
	err := ws.WriteJSON(errorMessage)
	if err != nil {
	}
}

func SendJSONResponse(ws *websocket.Conn, responseMessage fiber.Map) {
	err := ws.WriteJSON(responseMessage)
	if err != nil {
	}
}
