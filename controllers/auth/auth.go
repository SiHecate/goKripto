package controllers

import (
	model "gokripto/Model"
	"gokripto/database"
	"gokripto/helpers"
	helper "gokripto/helpers"
	"gokripto/repositories"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if data.Password != data.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Passwords do not match",
		})
	}

	passwordHash, err := helper.BcryptPassword(data.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to hash password",
		})
	}

	userRepo := repositories.NewUserRepository(database.Conn)
	user, err := userRepo.CreateUser(data, passwordHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	walletToken := helper.GenerateWalletToken()
	if err := userRepo.CreateWallet(walletToken, *user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create wallet",
		})
	}

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request data",
		})
	}

	user, err := model.GetUserByEmail(database.Conn, data.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password does not match",
		})
	}

	// JWT oluşturma işlemi
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * 2)).Time.Unix()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: expiresAt,
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not log in",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 2),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Success",
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
	issuer, err := helpers.GetIssuer(c)
	if err != nil {
		return err
	}

	user, err := model.GetUserByIssuer(database.Conn, issuer)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	walletAddress, err := model.GetWalletAddress(database.Conn, issuer)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Wallet not found",
		})
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
		WalletAddress: walletAddress,
		WalletBalance: user.Wallet.Balance,
	}
	return c.JSON(response)
}
