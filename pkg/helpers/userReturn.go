package helpers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetIssuerFromContext(c *fiber.Ctx) (string, error) {
	issuer, ok := c.Locals("issuer").(string)
	if !ok {
		return "", fmt.Errorf("issuer not found!")
	}
	return issuer, nil
}
