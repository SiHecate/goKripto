package controllers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestLogin(t *testing.T) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		var data map[string]string

		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request data",
			})
		}
		if _, ok := data["password"]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errorPassword": "Missing 'password' field",
			})
		}
		if _, ok := data["email"]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errorEmail": "Missing 'email' field",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Login successful",
		})
	})

	tests := []struct {
		name           string
		requestPayload string
		expected       wanted
	}{
		{
			name:           "Valid Registration (Login)",
			requestPayload: `{"email": "user123@user.com", "password": "user123"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing password)",
			requestPayload: `{"email": "user123@user.com"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorPassword",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing email)",
			requestPayload: `{"password": "user123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorEmail",
				},
			},
		},
		{
			name:           "Invalid Registration (Login)",
			requestPayload: `{"",}`,
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
			req := httptest.NewRequest("POST", "/login", strings.NewReader(test.requestPayload))
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

func BenchmarkLogin(b *testing.B) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		var data map[string]string

		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request data",
			})
		}
		if _, ok := data["password"]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errorPassword": "Missing 'password' field",
			})
		}
		if _, ok := data["email"]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errorEmail": "Missing 'email' field",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Login successful",
		})
	})

	tests := []struct {
		name           string
		requestPayload string
		expected       wanted
	}{
		{
			name:           "Valid Registration (Login)",
			requestPayload: `{"email": "user123@user.com", "password": "user123"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing password)",
			requestPayload: `{"email": "user123@user.com"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorPassword",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing email)",
			requestPayload: `{"password": "user123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorEmail",
				},
			},
		},
		{
			name:           "Invalid Registration (Login)",
			requestPayload: `{"",}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"error",
				},
			},
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			req := httptest.NewRequest("POST", "/login", strings.NewReader(test.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				b.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expected.StatusCode {
				b.Errorf("Expected status code %d, but got %d", test.expected.StatusCode, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				b.Fatal(err)
			}
			var response map[string]interface{}
			err = json.Unmarshal(body, &response)
			if err != nil {
				b.Fatal(err)
			}

			for _, key := range test.expected.expectedKeys {
				if _, ok := response[key]; !ok {
					b.Errorf("Expected JSON key '%s' not found in the response", key)
				}
			}
		})
	}
}
