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
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
		Wallet: model.Wallet{
			WalletAddress: walletToken,
		},
	}
	database.DB.Create(&user)

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
	database.DB.Where("wallet_address = ?", walletAddress).First(&user)
	return user
}

// CreateWallet creates a wallet for a user.
func CreateWallet(user model.User) error {
	wallet := model.Wallet{
		WalletAddress: user.Wallet.WalletAddress,
		UserID:        user.ID,
		Balance:       0,
	}
	database.DB.Create(&wallet)
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

const SecretKey = "secret"

// Login handles user login.
func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Parse the request body into the 'data' map.
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user model.User

	// Query the database to find the user by their email.
	database.DB.Where("email = ?", data["email"]).First(&user)

	// If no user is found, return a 404 response.
	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	// Compare the hashed password stored in the database with the provided password.
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "password not match",
		})
	}

	// Create a JWT token with the user's ID as the issuer and set an expiration time.
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
	})

	// Sign the JWT token with the secret key.
	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	// Set the JWT token as a cookie in the response.
	info := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 2),
		HTTPOnly: true,
	}

	c.Cookie(&info)

	// Return a JSON response indicating a successful login.
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

// Logout handles user logout by clearing the JWT cookie.
func Logout(c *fiber.Ctx) error {
	// Create a new cookie with an empty value and an expiration time in the past.
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	// Set the cookie in the response to clear the JWT cookie on the client side.
	c.Cookie(&cookie)

	// Return a JSON response indicating a successful logout.
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func User(c *fiber.Ctx) error {
	// Cookie and token
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

	// Return user info for table
	var user model.User
	if err := database.DB.Where("id = ?", claims.Issuer).Preload("Wallet").First(&user).Error; err != nil {
		return err
	}

	// responseData := map[string]interface{}{
	// 	"user_name":      user.Name,
	// 	"email":          user.Email,
	// 	"wallet_address": user.Wallet.WalletAddress,
	// 	"wallet_balance": user.Wallet.Balance,
	// }
	return c.JSON(user)
}
