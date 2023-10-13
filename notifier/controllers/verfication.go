package verification

import (
	model "/home/umut/goKripto/crypto/Model"
	"/home/umut/goKripto/crypto/database"
	helper "/home/umut/goKripto/crypto/helpers"

	"github.com/gofiber/fiber/v2"
)

// Diğer kodlar buraya gelir...

func Verfication(c *fiber.Ctx) error {
	issuer, err := helper.GetIssuer(c)
	if err != nil {
		return err
	}
	var data struct {
		Email       string `json:"email"`
		Verfication bool   `json:"verfication"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	user, err := model.GetUserByIssuer(database.Conn, issuer)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	type Verfication struct {
		Name  string
		Email string
	}

	response := Verfication{
		Name:  user.Name,
		Email: user.Email,
	}
	return c.JSON(response)
}
