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

func TestBalance(t *testing.T) {
	app := fiber.New()

	app.Get("/balance", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"walletID":      1,
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
				"walletID",
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
		t.Run(tt.description, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/balance", nil)

			res, err := app.Test(req)
			if err != nil {
				t.Errorf("Cannot test Fiber handler: %v", err)
				t.Fail()
			}

			if res.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, res.StatusCode)
			}

			if len(tt.expectedKeys) > 0 {
				var response map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Errorf("Cannot parse response body: %v", err)
				}

				for _, key := range tt.expectedKeys {
					if _, ok := response[key]; !ok {
						t.Errorf("Expected key '%s' not found in response", key)
					}
				}
			}
		})
	}
}

func BenchmarkBalance(b *testing.B) {
	app := fiber.New()

	app.Get("/balance", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"walletID":      1,
			"walletAddress": "denemeDENEME123456789",
			"userID":        1,
			"balance":       123123,
		}
		return c.JSON(response)
	})

	req := httptest.NewRequest("GET", "/balance", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := app.Test(req)
		if err != nil {
			b.Fatalf("Cannot test Fiber handler: %v", err)
		}

		if res.StatusCode != fiber.StatusOK {
			b.Fatalf("Expected status code %d, got %d", fiber.StatusOK, res.StatusCode)
		}

		var response map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			b.Fatalf("Cannot parse response body: %v", err)
		}

		expectedKeys := []string{
			"walletID",
			"walletAddress",
			"userID",
			"balance",
		}

		for _, key := range expectedKeys {
			if _, ok := response[key]; !ok {
				b.Fatalf("Expected key '%s' not found in response", key)
			}
		}

		res.Body.Close()
	}
}

func TestListAllCryptos(t *testing.T) {
	app := fiber.New()
	t.Parallel()

	app.Get("/cryptoList", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"cryptoID":   8,
			"cryptoname": "ETH",
			"price":      1626.79198,
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
				"cryptoID",
				"cryptoname",
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

			if res.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, res.StatusCode)
			}

			if len(tt.expectedKeys) > 0 {
				var response map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Errorf("Cannot parse response body: %v", err)
				}

				for _, key := range tt.expectedKeys {
					if _, ok := response[key]; !ok {
						t.Errorf("Expected key '%s' not found in response", key)
					}
				}
			}
		})
	}

}

func TestListCryptoWallet(t *testing.T) {
	app := fiber.New()
	t.Parallel()

	app.Get("/listcryptowallet", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"CryptoWalletID":   6,
			"walletAddress":    "%TE4i760HX6SNt5bMVURZ4dFR+%@@7w",
			"cryptoID":         7,
			"cryptoname":       "BTC",
			"cryptoTotalPrice": 3505818.36,
			"cryptoAmount":     132,
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
				"CryptoWalletID",
				"walletAddress",
				"cryptoID",
				"cryptoname",
				"cryptoTotalPrice",
				"cryptoAmount",
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

			if res.StatusCode != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, res.StatusCode)
			}

			if len(tt.expectedKeys) > 0 {
				var response map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Errorf("Cannot parse response body: %v", err)
				}

				for _, key := range tt.expectedKeys {
					if _, ok := response[key]; !ok {
						t.Errorf("Expected key '%s' not found in response", key)
					}
				}
			}
		})
	}

}

func BenchmarkListAllCryptos(b *testing.B) {
	app := fiber.New()

	app.Get("/cryptoList", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"cryptoID":   8,
			"cryptoname": "ETH",
			"price":      1626.79198,
		}

		return c.JSON(response)
	})

	req := httptest.NewRequest("GET", "/cryptoList", nil)

	for i := 0; i < b.N; i++ {
		res, err := app.Test(req)
		if err != nil {
			b.Fatalf("Cannot test Fiber handler: %v", err)
		}

		if res.StatusCode != 200 {
			b.Fatalf("Expected status code 200, got %d", res.StatusCode)
		}
	}
}

func BenchmarkListCryptoWallet(b *testing.B) {
	app := fiber.New()

	app.Get("/listcryptowallet", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"CryptoWalletID":   6,
			"walletAddress":    "%TE4i760HX6SNt5bMVURZ4dFR+%@@7w",
			"cryptoID":         7,
			"cryptoname":       "BTC",
			"cryptoTotalPrice": 3505818.36,
			"cryptoAmount":     132,
		}

		return c.JSON(response)
	})

	req := httptest.NewRequest("GET", "/listcryptowallet", nil)

	for i := 0; i < b.N; i++ {
		res, err := app.Test(req)
		if err != nil {
			b.Fatalf("Cannot test Fiber handler: %v", err)
		}

		if res.StatusCode != 200 {
			b.Fatalf("Expected status code 200, got %d", res.StatusCode)
		}
	}
}

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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/cryptoBuy", strings.NewReader(test.requestPayload))
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
