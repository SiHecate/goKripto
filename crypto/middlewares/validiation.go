package middlewares

import (
	model "cryptoApp/Model"
	"cryptoApp/database"
	"cryptoApp/helpers"

	"github.com/gofiber/fiber/v2"
)

func Validiation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		issuer, err := helpers.GetIssuer(c)
		if err != nil {
			return err
		}

		var user model.User
		if err := database.Conn.Where("id = ?", issuer).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "User not found",
				"issuer":  issuer,
			})
		}

		if user.Verfication != true {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Account not verified",
			})
		}

		return c.Next()
	}
}
