package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Logout handles user logout by clearing the JWT cookie.
func Logout(c *fiber.Ctx) error {
	// Create a new cookie with an empty value and an expiration time in the past.
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	// Set the cookie in the response to clear the JWT cookie on the client side.
	c.Cookie(&cookie)

	// Return a JSON response indicating a successful logout.
	return c.JSON(fiber.Map{
		"message": "success",
	})
}
