package controllers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestBuyCryptos(t *testing.T) {
	app := fiber.New()
	app.Post("/cryptoBuy", func(c *fiber.Ctx) error {
		return nil
	})
	req := httptest.NewRequest("POST", "/cryptoBuy", strings.NewReader(`{
		"cryptoName": "Bitcoin",
		"amountToBuy": 1.0
	}`))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status code %d, but got %d", fiber.StatusOK, resp.StatusCode)
	}
}

func BenchmarkBuyCryptos(b *testing.B) {
	app := fiber.New()
	app.Post("/cryptoBuy", func(c *fiber.Ctx) error {
		return nil
	})
	req := httptest.NewRequest("POST", "/cryptoBuy", strings.NewReader(`{
		"cryptoName": "Bitcoin",
		"amountToBuy": 1.0
	}`))
	req.Header.Set("Content-Type", "application/json")
	for i := 0; i < b.N; i++ {
		resp, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != fiber.StatusOK {
			b.Errorf("Expected status code %d, but got %d", fiber.StatusOK, resp.StatusCode)
		}
	}
}
