package helpers

import (
	model "cryptoApp/Model"
	"errors"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func BcryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ComparePasswords(password string) (bool, error) {
	var user model.User
	passwordHash := string(user.Password)
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return false, err
}

func GenerateWalletToken() string {

	const chars = "1234567890"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := r.Intn(12) + 26
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}
