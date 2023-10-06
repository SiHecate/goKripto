package controllers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Done Get
func TestBalance(t *testing.T) {
	app := fiber.New()

	app.Get("/balance", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"wallet_address": "denemeDENEME123456789",
			"username":       "deneme",
			"balance":        123123,
		}
		return c.JSON(response)
	})

	tests := []struct {
		description  string
		statusCode   int
		expectedKeys []string
	}{
		{
			description: "get HTTP status 200",
			statusCode:  200,
			expectedKeys: []string{
				"wallet_address",
				"username",
				"balance",
			},
		},
		{
			description:  "get HTTP status 200 without expected keys",
			statusCode:   200,
			expectedKeys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/balance", nil)

			res, err := app.Test(req)
			if err != nil {
				t.Errorf("Cannot test Fiber handler: %v", err)
				t.Fail()
			}

			assert.Equal(t, tt.statusCode, res.StatusCode, "Expected status code %d, got %d", tt.statusCode, res.StatusCode)

			if len(tt.expectedKeys) > 0 {
				var response map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Errorf("Cannot parse response body: %v", err)
				}

				for _, key := range tt.expectedKeys {
					assert.Contains(t, response, key, "Expected key '%s' not found in response", key)
				}
			}
		})
	}
}

// Done Get
func BenchmarkBalance(b *testing.B) {
	app := fiber.New()

	app.Get("/balance", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"walletAddress": "denemeDENEME123456789",
			"userID":        1,
			"balance":       123123,
		}
		return c.JSON(response)
	})

	tests := []struct {
		description  string
		statusCode   int
		expectedKeys []string
	}{
		{
			description: "get HTTP status 200",
			statusCode:  200,
			expectedKeys: []string{
				"walletAddress",
				"userID",
				"balance",
			},
		},
		{
			description:  "get HTTP status 200 without expected keys",
			statusCode:   200,
			expectedKeys: []string{},
		},
	}

	for _, tt := range tests {
		b.Run(tt.description, func(b *testing.B) {
			req := httptest.NewRequest("GET", "/balance", nil)

			res, err := app.Test(req)
			if err != nil {
				b.Errorf("Cannot test Fiber handler: %v", err)
				b.Fail()
			}

			assert.Equal(b, tt.statusCode, res.StatusCode, "Expected status code %d, got %d", tt.statusCode, res.StatusCode)

			if len(tt.expectedKeys) > 0 {
				var response map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					b.Errorf("Cannot parse response body: %v", err)
				}

				for _, key := range tt.expectedKeys {
					assert.Contains(b, response, key, "Expected key '%s' not found in response", key)
				}
			}
		})
	}
}

// Done Get
func TestListAllCryptos(t *testing.T) {
	app := fiber.New()
	t.Parallel()

	app.Get("/cryptoList", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"crypto_symbol": "ETH",
			"crypto_name":   "Ethereum",
			"price":         1623.9262526599507,
		}

		return c.JSON(response)
	})

	tests := []struct {
		description  string
		statusCode   int
		expectedKeys []string
	}{
		{
			description: "get HTTP status 200",
			statusCode:  200,
			expectedKeys: []string{
				"crypto_symbol",
				"crypto_name",
				"price",
			},
		},
		{
			description:  "get HTTP status 200 without expected keys",
			statusCode:   200,
			expectedKeys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/cryptoList", nil)

			res, err := app.Test(req)
			if err != nil {
				t.Errorf("Cannot test Fiber handler: %v", err)
				t.Fail()
			}

			assert.Equal(t, tt.statusCode, res.StatusCode, "Expected status code %d, got %d", tt.statusCode, res.StatusCode)

			if len(tt.expectedKeys) > 0 {
				var response map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Errorf("Cannot parse response body: %v", err)
				}

				for _, key := range tt.expectedKeys {
					assert.Contains(t, response, key, "Expected key '%s' not found in response", key)
				}
			}
		})
	}

}

// Done Get
func BenchmarkListAllCryptos(b *testing.B) {
	app := fiber.New()

	app.Get("/cryptoList", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"crypto_symbol": "ETH",
			"crypto_name":   "Ethereum",
			"price":         1623.9262526599507,
		}

		return c.JSON(response)
	})

	tests := []struct {
		description  string
		statusCode   int
		expectedKeys []string
	}{
		{
			description: "get HTTP status 200",
			statusCode:  200,
			expectedKeys: []string{
				"crypto_symbol",
				"crypto_name",
				"price",
			},
		},
		{
			description:  "get HTTP status 200 without expected keys",
			statusCode:   200,
			expectedKeys: []string{},
		},
	}

	for _, tt := range tests {
		b.Run(tt.description, func(b *testing.B) {
			req := httptest.NewRequest("GET", "/cryptoList", nil)

			res, err := app.Test(req)
			if err != nil {
				b.Errorf("Cannot test Fiber handler: %v", err)
				b.Fail()
			}

			assert.Equal(b, tt.statusCode, res.StatusCode, "Expected status code %d, got %d", tt.statusCode, res.StatusCode)

			if len(tt.expectedKeys) > 0 {
				var response map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					b.Errorf("Cannot parse response body: %v", err)
				}

				for _, key := range tt.expectedKeys {
					assert.Contains(b, response, key, "Expected key '%s' not found in response", key)
				}
			}
		})
	}

}

// Done Get
func TestListCryptoWallet(t *testing.T) {
	app := fiber.New()
	t.Parallel()

	app.Get("/listcryptowallet", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"wallet_address":     "E4i760HX6SNt5bMVURZ4dFR+%@@7w",
			"crypto_name":        "Bitcoin",
			"crypto_amount":      7.36,
			"crypto_total_price": 186936.52406785818,
		}

		return c.JSON(response)
	})

	tests := []struct {
		description  string
		statusCode   int
		expectedKeys []string
	}{
		{
			description: "get HTTP status 200",
			statusCode:  200,
			expectedKeys: []string{
				"wallet_address",
				"crypto_name",
				"crypto_amount",
				"crypto_total_price",
			},
		},
		{
			description:  "get HTTP status 200 without expected keys",
			statusCode:   200,
			expectedKeys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/listcryptowallet", nil)

			res, err := app.Test(req)
			if err != nil {
				t.Errorf("Cannot test Fiber handler: %v", err)
				t.Fail()
			}

			assert.Equal(t, tt.statusCode, res.StatusCode, "Expected status code %d, got %d", tt.statusCode, res.StatusCode)

			if len(tt.expectedKeys) > 0 {
				var response map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Errorf("Cannot parse response body: %v", err)
				}

				for _, key := range tt.expectedKeys {
					assert.Contains(t, response, key, "Expected key '%s' not found in response", key)
				}
			}
		})
	}

}

// Done Get
func BenchmarkListCryptoWallet(b *testing.B) {
	app := fiber.New()

	app.Get("/listcryptowallet", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"wallet_address":     "E4i760HX6SNt5bMVURZ4dFR+%@@7w",
			"crypto_name":        "Bitcoin",
			"crypto_amount":      7.36,
			"crypto_total_price": 186936.52406785818,
		}

		return c.JSON(response)
	})

	tests := []struct {
		description  string
		statusCode   int
		expectedKeys []string
	}{
		{
			description: "get HTTP status 200",
			statusCode:  200,
			expectedKeys: []string{
				"wallet_address",
				"crypto_name",
				"crypto_amount",
				"crypto_total_price",
			},
		},
		{
			description:  "get HTTP status 200 without expected keys",
			statusCode:   200,
			expectedKeys: []string{},
		},
	}

	for _, tt := range tests {
		b.Run(tt.description, func(b *testing.B) {
			req := httptest.NewRequest("GET", "/listcryptowallet", nil)

			res, err := app.Test(req)
			if err != nil {
				b.Errorf("Cannot test Fiber handler: %v", err)
				b.Fail()
			}

			assert.Equal(b, tt.statusCode, res.StatusCode, "Expected status code %d, got %d", tt.statusCode, res.StatusCode)

			if len(tt.expectedKeys) > 0 {
				var response map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					b.Errorf("Cannot parse response body: %v", err)
				}

				for _, key := range tt.expectedKeys {
					assert.Contains(b, response, key, "Expected key '%s' not found in response", key)
				}
			}
		})
	}

}

// Done Post
func TestBuyCrypto(t *testing.T) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/cryptoBuy", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}

		if data["cryptoName"] == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"nameError": "Crypto name does not exist",
			})
		}

		amountStr := data["amountToBuy"]
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
			name:           "Valid Registration Buy Crypto",
			requestPayload: `{"cryptoName": "btc", "amountToBuy": "123"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Name Error)",
			requestPayload: `{"amountToBuy": "123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"nameError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Error)",
			requestPayload: `{"cryptoName": "BTC", "amountToBuy": ""}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"amountError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Type Error)",
			requestPayload: `{"cryptoName": "BTC", "amountToBuy": "asd123"}`,
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
			req := httptest.NewRequest("POST", "/cryptoBuy", strings.NewReader(test.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, test.expected.StatusCode, resp.StatusCode)

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
				assert.Contains(t, response, key, "Expected JSON key '%s' not found in the response", key)
			}
		})
	}
}

// Done Post
func BenchmarkBuyCrypto(b *testing.B) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/cryptoBuy", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}

		if data["cryptoName"] == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"nameError": "Crypto name does not exists",
			})
		}

		amountStr := data["amountToBuy"]
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
			name:           "Valid Registration Buy Crypto",
			requestPayload: `{"cryptoName": "btc", "amountToBuy": "123"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Name Error)",
			requestPayload: `{"amountToBuy": "123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"nameError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Error)",
			requestPayload: `{"cryptoName": "BTC", "amountToBuy": ""}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"amountError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Type Error)",
			requestPayload: `{"cryptoName": "BTC", "amountToBuy": "asd123"}`,
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
			req := httptest.NewRequest("POST", "/cryptoBuy", strings.NewReader(benchmark.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				b.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(b, benchmark.expected.StatusCode, resp.StatusCode)

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
				assert.Contains(b, response, key, "Expected JSON key '%s' not found in the response", key)
			}
		})
	}
}

// Done Post
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
				"nameError": "Crypto name does not exist",
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

			assert.Equal(t, test.expected.StatusCode, resp.StatusCode)

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
				assert.Contains(t, response, key, "Expected JSON key '%s' not found in the response", key)
			}
		})
	}
}

// Done Post
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

			assert.Equal(b, benchmark.expected.StatusCode, resp.StatusCode)

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
				assert.Contains(b, response, key, "Expected JSON key '%s' not found in the response", key)
			}
		})
	}
}

// Done Post
func TestAddBalance(t *testing.T) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/addBalance", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}

		additiveBalanceStr := data["addBalance"]
		if additiveBalanceStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"amountError": "Amount does not exist",
			})
		}

		additiveBalance, err := strconv.ParseFloat(additiveBalanceStr, 64)
		if err != nil || additiveBalance <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"amountTypeError": "Invalid type of amount",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Adding balance successful",
		})
	})

	tests := []struct {
		name           string
		requestPayload string
		expected       wanted
	}{
		{
			name:           "Valid Registration (Add Balance)",
			requestPayload: `{"addBalance": "1000000"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Error)",
			requestPayload: `{"addBalance": ""}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"amountError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Type Error)",
			requestPayload: `{"addBalance": "asd123"}`,
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
			req := httptest.NewRequest("POST", "/addBalance", strings.NewReader(test.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(t, test.expected.StatusCode, resp.StatusCode)

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
				assert.Contains(t, response, key, "Expected JSON key '%s' not found in the response", key)
			}
		})
	}
}

// Done Post
func BenchmarkAddBalance(b *testing.B) {
	type wanted struct {
		StatusCode   int
		expectedKeys []string
	}

	app := fiber.New()

	app.Post("/addBalance", func(c *fiber.Ctx) error {
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return err
		}

		additiveBalanceStr := data["addBalance"]
		if additiveBalanceStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"amountError": "Amount does not exist",
			})
		}

		additiveBalance, err := strconv.ParseFloat(additiveBalanceStr, 64)
		if err != nil || additiveBalance <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"amountTypeError": "Invalid type of amount",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Adding balance successful",
		})
	})

	tests := []struct {
		name           string
		requestPayload string
		expected       wanted
	}{
		{
			name:           "Valid Registration (Add Balance)",
			requestPayload: `{"addBalance": "1000000"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Error)",
			requestPayload: `{"addBalance": ""}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"amountError",
				},
			},
		},
		{
			name:           "Unvalid Registration (Amount Type Error)",
			requestPayload: `{"addBalance": "asd123"}`,
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
			req := httptest.NewRequest("POST", "/addBalance", strings.NewReader(benchmark.requestPayload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				b.Fatal(err)
			}
			defer resp.Body.Close()

			assert.Equal(b, benchmark.expected.StatusCode, resp.StatusCode)

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
				assert.Contains(b, response, key, "Expected JSON key '%s' not found in the response", key)
			}
		})
	}
}

func TestTransactionListCrypto(t *testing.T) {
	app := fiber.New()

	app.Get("/TransactionListCrypto", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"walletAddress": "denemeDENEME123456789",
			"userID":        1,
			"balance":       123123,
		}
		return c.JSON(response)
	})
}

func TestTransactionListBalance(t *testing.T) {

}
