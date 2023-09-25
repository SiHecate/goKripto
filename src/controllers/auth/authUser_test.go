package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestUser(t *testing.T) {
	app := fiber.New()

	app.Get("/user", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"userID":        2,
			"username":      "umut",
			"email":         "umut@admin.com",
			"walletAddress": "%TE4i760HX6SNt5bMVURZ4dFR+%@@7w",
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
				"userID",
				"username",
				"email",
				"walletAddress",
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
			req := httptest.NewRequest("GET", "/user", nil)

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
