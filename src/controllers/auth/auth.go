package controllers

import (
	model "gokripto/Model"
	"gokripto/database"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	walletToken := generateWalletToken()

	user := model.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
		Wallet: model.Wallet{
			WalletAddress: walletToken,
		},
	}
	database.DB.Create(&user)

	createdUser, err := GetUserByWalletAddress(database.DB, walletToken)
	if err != nil {
		// Handle the error here
		return err
	}

	CreateWallet(*createdUser)

	return c.JSON(user)
}

func GetUserByWalletAddress(db *gorm.DB, walletAddress string) (*model.User, error) {
	var user model.User
	if err := db.Where("wallets.wallet_address = ?", walletAddress).
		Joins("JOIN wallets ON wallets.user_id = users.id").
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateWallet(user model.User) error {
	wallet := model.Wallet{
		WalletAddress: user.Wallet.WalletAddress,
		UserID:        user.ID,
		Balance:       0,
	}
	database.DB.Create(&wallet)
	return nil
}

func generateWalletToken() string {

	const chars = "1234567890"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := r.Intn(10) + 26
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

const SecretKey = "secret"

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user model.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "password not match",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	info := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 2),
		HTTPOnly: true,
	}

	c.Cookie(&info)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)

	var user model.User
	if err := database.DB.Where("id = ?", claims.Issuer).Preload("Wallet").First(&user).Error; err != nil {
		return err
	}

	type UserResponse struct {
		Name          string
		Email         string
		WalletAddress string
		WalletBalance float64
	}

	response := UserResponse{
		Name:          user.Name,
		Email:         user.Email,
		WalletAddress: user.Wallet.WalletAddress,
		WalletBalance: user.Wallet.Balance,
	}
	return c.JSON(response)
}
