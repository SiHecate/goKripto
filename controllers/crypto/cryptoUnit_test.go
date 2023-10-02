package controllers

import (
	model "gokripto/Model"
	"gokripto/database"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func TestAddBalanceCrypto(t *testing.T) {
	// Fiber uygulamasını başlat
	app := fiber.New()

	// fiber.Ctx nesnesini oluştur
	c := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(c)

	// Test için uygun verileri oluştur
	validData := `{"add_balance": 100.0}`
	c.Request().SetBodyString(validData)

	// Sahte uygulama oluştur
	mockApp := new(MockControllersApp)

	mockApp.On("AddBalanceCrypto", mock.Anything).Return(nil)

	mockApp.On("Model", &model.Wallet{}).Return(database.Conn)
	mockApp.On("Update", "balance", float64(200.0)).Return(nil)

	err := mockApp.AddBalanceCrypto(c)

	assert.NoError(t, err)

	// Sahte çağrıların beklenen şekilde çalıştığını doğrula
	mockApp.AssertExpectations(t)
}
