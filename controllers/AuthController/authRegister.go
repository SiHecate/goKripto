package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	walletToken := generateWalletToken()
	user := model.User{
		Name:          data["name"],
		Email:         data["email"],
		Password:      password,
		WalletAddress: walletToken,
	}
	Database.DB.Create(&user)

	createdUser := GetUserWalletAddress(walletToken)
	CreateWallet(createdUser)

	return c.JSON(user)
}

func GetUserWalletAddress(walletAddress string) model.User {
	var user model.User
	Database.DB.Where("wallet_address = ?", walletAddress).First(&user)
	return user
}

func CreateWallet(user model.User) error {
	wallet := model.Wallet{
		WalletAddress: user.WalletAddress,
		UserID:        user.Id,
		Balance:       0,
	}
	Database.DB.Create(&wallet)
	return nil
}

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
