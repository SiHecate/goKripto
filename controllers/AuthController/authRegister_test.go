package controllers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegister(t *testing.T) {

	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/register", func(c *fiber.Ctx) error {
		var data map[string]string

		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request data",
			})
		}
		if _, ok := data["password"]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing 'password' field",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Registration successful",
		})
	})

	tests := []struct {
		name           string
		requestPayload string
		expected       wanted
	}{
		{
			name:           "Valid Registration",
			requestPayload: `{"name": "user123", "email": "user123@user.com", "password": "user123"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing Password)",
			requestPayload: `{"name": "user123", "email": "user123@user.com"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"error",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/register", strings.NewReader(test.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expected.StatusCode {
				t.Errorf("Expected status code %d, but got %d", test.expected.StatusCode, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			var response map[string]interface{}
			err = json.Unmarshal(body, &response)
			if err != nil {
				t.Fatal(err)
			}

			for _, key := range test.expected.expectedKeys {
				if _, ok := response[key]; !ok {
					t.Errorf("Expected JSON key '%s' not found in the response", key)
				}
			}
		})
	}
}
