package verification

import (
	"Notifier/database"
	"Notifier/model"

	"github.com/gofiber/fiber/v2"
)

func VerificationTable()

func Verification(c *fiber.Ctx) error {
	var data struct {
		Email        string `json:"email"`
		Verification string `json:"verification"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Geçersiz istek verisi",
		})
	}

	email := data.Email
	verificationCode := data.Verification

	var emailVerification model.Verfication
	if err := database.Conn.Where("email = ?", email).First(&emailVerification).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Mail doğrulaması başarısız",
			"error":   err.Error(),
		})
	}

	if emailVerification.Verfication != verificationCode {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Kod doğrulaması başarısız",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Doğrulama başarılı",
	})
}

func PrintVerificationTable(c *fiber.Ctx) error {
	var verifications []model.Verfication

	if err := database.Conn.Find(&verifications).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Tabloya erişim başarısız",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(verifications)
}
