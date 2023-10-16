package verification

import (
	"Notifier/database"
	"Notifier/model"

	"github.com/gofiber/fiber/v2"
)

func Verification(c *fiber.Ctx) error {
	var data struct {
		Email        string `json:"email"`
		Verification bool   `json:"verification"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	email := data.Email
	verificationCode := data.Verification

	// Check verification
	verificationQuery := database.Conn.Model(&model.Verfication{}).Where("email = ? OR verification = ?", email, verificationCode)
	if err := verificationQuery.Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Verification failed",
			"error":   err.Error(),
		})
	}

	// Update 'verification' field
	userUpdateQuery := database.Conn.Model(&model.User{}).Where("email = ?", email).Update("verification", true)
	if updateErr := userUpdateQuery.Error; updateErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Error updating verification",
			"error":   updateErr.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
