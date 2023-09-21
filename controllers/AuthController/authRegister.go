package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register handles user registration.
func Register(c *fiber.Ctx) error {
	var data map[string]string

	// Parse the request body into the 'data' map.
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Hash the user's password.
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	// Generate a wallet token.
	walletToken := generateWalletToken()

	// Create a new user record in the database.
	user := model.User{
		Name:          data["name"],
		Email:         data["email"],
		Password:      password,
		WalletAddress: walletToken,
	}
	Database.DB.Create(&user)

	// Retrieve the created user with wallet address.
	createdUser := GetUserWalletAddress(walletToken)

	// Create a wallet for the user.
	CreateWallet(createdUser)

	// Return the JSON representation of the created user.
	return c.JSON(user)
}

// GetUserWalletAddress retrieves a user by their wallet address.
func GetUserWalletAddress(walletAddress string) model.User {
	var user model.User
	Database.DB.Where("wallet_address = ?", walletAddress).First(&user)
	return user
}

// CreateWallet creates a wallet for a user.
func CreateWallet(user model.User) error {
	wallet := model.Wallet{
		WalletAddress: user.WalletAddress,
		UserID:        user.Id,
		Balance:       0,
	}
	Database.DB.Create(&wallet)
	return nil
}

// generateWalletToken generates a random wallet token.
func generateWalletToken() string {

	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789%+@"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := r.Intn(10) + 26
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
