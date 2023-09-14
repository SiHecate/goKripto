package controllers

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestSellCryptos(t *testing.T) {
	app := fiber.New()

	app.Post("/CryptoSell", func(c *fiber.Ctx) error {
		return nil
	})
	req := httptest.NewRequest("POST", "/cryptoSell", strings.NewReader(`{
		"cryptoName": "Bitcoin",
		"amountToSell": 11.00
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

	var response struct{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkSellCryptos(b *testing.B) {
	app := fiber.New()
	app.Post("/cryptoSell", func(c *fiber.Ctx) error {
		return nil
	})
	req := httptest.NewRequest("POST", "/cryptoSell", strings.NewReader(`{
		"cryptoName": "Bitcoin",
		"amountToBuy": 11.23
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
