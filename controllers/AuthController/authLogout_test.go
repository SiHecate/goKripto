package controllers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestLogout(t *testing.T) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/logout", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}

		if data["Name"] != "jwt" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"nameError": "Name must be jwt",
			})
		}

		if data["Value"] != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"valueError": "Value must be empty",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logout successful",
		})
	})

	tests := []struct {
		name           string
		requestPayload string
		expected       wanted
	}{
		{
			name:           "Valid Registration (Logout)",
			requestPayload: `{"Name": "jwt", "Value": ""}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Value Error)",
			requestPayload: `{"Name": "jwt", "Value": "123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"valueError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Name Error)",
			requestPayload: `{"Name": "asdf123", "Value": "123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"nameError",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/logout", strings.NewReader(test.requestPayload))
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

func BenchmarkLogout(b *testing.B) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/logout", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}

		if data["Name"] != "jwt" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"nameError": "Name must be jwt",
			})
		}

		if data["Value"] != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"valueError": "Value must be empty",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logout successful",
		})
	})

	tests := []struct {
		name           string
		requestPayload string
		expected       wanted
	}{
		{
			name:           "Valid Registration (Logout)",
			requestPayload: `{"Name": "jwt", "Value": ""}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Value Error)",
			requestPayload: `{"Name": "jwt", "Value": "123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"valueError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Name Error)",
			requestPayload: `{"Name": "asdf123", "Value": "123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"nameError",
				},
			},
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			req := httptest.NewRequest("POST", "/logout", strings.NewReader(test.requestPayload))
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
