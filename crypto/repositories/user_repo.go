package repositories

import (
	model "cryptoApp/Model"
	"cryptoApp/database"
	"cryptoApp/helpers"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (repo *UserRepository) CreateWallet(walletToken string, user model.User) error {
	wallet := model.Wallet{
		WalletAddress: walletToken,
		UserID:        user.ID,
		Balance:       0,
	}

	database.Conn.Create(&wallet)

	return nil
}

func (repo *UserRepository) CreateUser(data struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}, passwordHash string) (*model.User, error) {
	passwordBytes := []byte(passwordHash)

	user := model.User{
		Name:     data.Name,
		Email:    data.Email,
		Password: passwordBytes,
	}

	verfication_code := helpers.GenerateCode()

	verfication := model.Verfication{
		Email:       data.Email,
		Verfication: verfication_code,
	}

	if err := repo.db.Create(&user).Error; err != nil {
		return nil, err
	}
	if err := repo.db.Create(&verfication).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
