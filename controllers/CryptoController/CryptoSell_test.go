package controllers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSellCrypto(t *testing.T) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/cryptoSell", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}

		if data["cryptoName"] == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"nameError": "Crypto name does not exists",
			})
		}

		amountStr := data["amountToSell"]
		if amountStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"amountError": "Amount does not exist",
			})
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"amountTypeError": "Invalid type of amount",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Trade successful",
		})
	})

	tests := []struct {
		name           string
		requestPayload string
		expected       wanted
	}{
		{
			name:           "Valid Registration (Sell Crypto)",
			requestPayload: `{"cryptoName": "btc", "amountToSell": "123"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Name Error)",
			requestPayload: `{"amountToSell": "123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"nameError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Error)",
			requestPayload: `{"cryptoName": "BTC", "amountToSell": ""}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"amountError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Type Error)",
			requestPayload: `{"cryptoName": "BTC", "amountToSell": "asd123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"amountTypeError",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/cryptoSell", strings.NewReader(test.requestPayload))
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

func BenchmarkSellCrypto(b *testing.B) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/cryptoSell", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}

		if data["cryptoName"] == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"nameError": "Crypto name does not exists",
			})
		}

		amountStr := data["amountToSell"]
		if amountStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"amountError": "Amount does not exist",
			})
		}

		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil || amount <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"amountTypeError": "Invalid type of amount",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Trade successful",
		})
	})

	tests := []struct {
		name           string
		requestPayload string
		expected       wanted
	}{
		{
			name:           "Valid Registration (Sell Crypto)",
			requestPayload: `{"cryptoName": "btc", "amountToSell": "123"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Name Error)",
			requestPayload: `{"amountToSell": "123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"nameError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Error)",
			requestPayload: `{"cryptoName": "BTC", "amountToSell": ""}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"amountError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Type Error)",
			requestPayload: `{"cryptoName": "BTC", "amountToSell": "asd123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"amountTypeError",
				},
			},
		},
	}

	for _, benchmark := range tests {
		b.Run(benchmark.name, func(b *testing.B) {
			req := httptest.NewRequest("POST", "/cryptoSell", strings.NewReader(benchmark.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				b.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != benchmark.expected.StatusCode {
				b.Errorf("Expected status code %d, but got %d", benchmark.expected.StatusCode, resp.StatusCode)
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

			for _, key := range benchmark.expected.expectedKeys {
				if _, ok := response[key]; !ok {
					b.Errorf("Expected JSON key '%s' not found in the response", key)
				}
			}
		})
	}
}
