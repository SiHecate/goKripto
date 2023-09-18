package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

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
