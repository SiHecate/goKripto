package repositories

import (
	model "gokripto/Model"
	"gokripto/database"

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

	if err := repo.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepository) GetWalletAddress(issuer string) (string, error) {
	modelWallet := model.Wallet{}

	if err := database.Conn.Where("user_id = ?", issuer).First(&modelWallet).Error; err != nil {
		return "", err
	}

	WalletAddress := modelWallet.WalletAddress
	return WalletAddress, nil
}
