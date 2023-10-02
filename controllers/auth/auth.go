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

type RegisterResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Name          string  `json:"name"`
	Email         string  `json:"email"`
	WalletAddress string  `json:"wallet_address"`
	WalletBalance float64 `json:"wallet_balance"`
}

type Status400Response struct {
	Message string `json:"StatusBadRequest"`
}

type Status401Response struct {
	Message string `json:"StatusUnauthorized"`
}

type Status404Response struct {
	Message string `json:"StatusNotFound"`
}

// Register
// @Summary Register user
// @Description Register func for new user
// @Tags User
// @Accept json
// @Produce json
// @Param name body string true "user name for register example: Umutcan "
// @Param email body string true "email address for register example"
// @Param password body string true "password for register example: umutcan123 or umutcan number not required "
// @Param confrim_pasword body string true "confirm password for register example: umutcan123 or umutcan"
// @Success 200 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Router /user/register [post]

func Register(c *fiber.Ctx) error {
	var data struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"` // Update the field name to match the struct definition
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

// Login
// @Summary Authenticate and log in a user
// @Description Authenticates a user by email and password and generates a JWT token.
// @Tags User
// @Accept json
// @Produce json
// @Param  email body string true "user email for log in example: Umutcan@example.com"
// @Param  password body string ture "password for log in example: umutcan123 or umutcan number not required "
// @Success 200 {object} LoginResponse
// @Failure 400 {object} Status400Response
// @Failure 401 {object} Status401Response
// @Router /user/login [post]
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

// Logout
// @Summary Log out the user
// @Description Logs out the authenticated user by clearing the JWT token cookie.
// @Tags User
// @Produce json
// @Success 200 {object} SuccessResponse
// @Router /user/logout [post]
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

// User
// @Summary Get user information
// @Description Get the user's information including name, email, wallet address, and wallet balance.
// @Tags User
// @Produce json
// @Success 200 {object} UserResponse
// @Failure 404 {object} Status404Response
// @Router /user/user [get]
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
