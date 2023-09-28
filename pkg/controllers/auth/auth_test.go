package controllers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Done Post
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

	for _, benchmark := range tests {
		b.Run(benchmark.name, func(b *testing.B) {
			req := httptest.NewRequest("POST", "/login", strings.NewReader(benchmark.requestPayload))
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

	for _, benchmark := range tests {
		b.Run(benchmark.name, func(b *testing.B) {
			req := httptest.NewRequest("POST", "/logout", strings.NewReader(benchmark.requestPayload))
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

// Done Get
func TestUser(t *testing.T) {
	app := fiber.New()

	app.Get("/user", func(c *fiber.Ctx) error {
		response := map[string]interface{}{
			"Name":          2,
			"Email":         "umut",
			"WalletAddress": "TE4i760HX6SNt5bMVURZ4dFR+%@@7w",
			"WalletBalance": "2333303.1321",
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
				"Name",
				"Email",
				"WalletAddress",
				"WalletBalance",
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

// Done Post
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
				"errorPassword": "Missing 'password' field",
			})
		}
		if _, ok := data["name"]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errorName": "Missing 'name' field",
			})
		}
		if _, ok := data["email"]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errorEmail": "Missing 'email' field",
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
			name:           "Valid Registration (Register)",
			requestPayload: `{"name": "user123", "email": "user123@user.com", "password": "user123"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Register)",
			requestPayload: `{""}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"error",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing Password)",
			requestPayload: `{"name": "user123", "email": "user123@user.com"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorPassword",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing Username)",
			requestPayload: `{"email": "user123@user.com", "password": "user123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorName",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing Email)",
			requestPayload: `{"name": "user123", "password": "user123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorEmail",
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
func BenchmarkRegister(b *testing.B) {

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
				"errorPassword": "Missing 'password' field",
			})
		}
		if _, ok := data["name"]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errorName": "Missing 'name' field",
			})
		}
		if _, ok := data["email"]; !ok {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errorEmail": "Missing 'email' field",
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
			name:           "Valid Registration (Register)",
			requestPayload: `{"name": "user123", "email": "user123@user.com", "password": "user123"}`,
			expected: wanted{
				StatusCode: 200,
				expectedKeys: []string{
					"message",
				},
			},
		},
		{
			name:           "Unvalid Registration (Register)",
			requestPayload: `{""}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"error",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing Password)",
			requestPayload: `{"name": "user123", "email": "user123@user.com"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorPassword",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing Username)",
			requestPayload: `{"email": "user123@user.com", "password": "user123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorName",
				},
			},
		},
		{
			name:           "Invalid Registration (Missing Email)",
			requestPayload: `{"name": "user123", "password": "user123"}`,
			expected: wanted{
				StatusCode: 400,
				expectedKeys: []string{
					"errorEmail",
				},
			},
		},
	}

	for _, benchmark := range tests {
		b.Run(benchmark.name, func(b *testing.B) {
			req := httptest.NewRequest("POST", "/register", strings.NewReader(benchmark.requestPayload))
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
