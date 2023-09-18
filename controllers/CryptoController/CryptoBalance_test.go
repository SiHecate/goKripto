package controllers

import (
	"encoding/json"
	"net/http/httptest"
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
