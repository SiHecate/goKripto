package controllers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
)

type MockControllersApp struct {
	mock.Mock
}

func (m *MockControllersApp) AddAllCryptoData(ws *websocket.Conn) error {
	args := m.Called(ws)
	return args.Error(0)
}

func (m *MockControllersApp) ListAllCryptos(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockControllersApp) AccountBalance(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockControllersApp) AddBalanceCrypto(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockControllersApp) BuyCryptos(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockControllersApp) SellCryptos(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockControllersApp) TransactionListBalance(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockControllersApp) TransactionListCrypto(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockControllersApp) ListCryptoWallet(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}
