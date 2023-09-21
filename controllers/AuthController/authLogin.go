package controllers

import (
	"gokripto/Database"
	model "gokripto/Model"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

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
	Database.DB.Where("email = ?", data["email"]).First(&user)

	// If no user is found, return a 404 response.
	if user.Id == 0 {
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
		Issuer:    strconv.Itoa(int(user.Id)),
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
