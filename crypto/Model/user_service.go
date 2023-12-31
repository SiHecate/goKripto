package model

import (
	"fmt"

	"gorm.io/gorm"
)

func GetUserByWalletAddress(db *gorm.DB, walletAddress string) (*User, error) {
	var user User
	if err := db.
		Where("wallets.wallet_address = ?", walletAddress).
		Joins("JOIN wallets ON wallets.user_id = users.id").
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByIssuer(db *gorm.DB, issuer string) (*User, error) {
	var user User
	if err := db.Where("id = ?", issuer).Preload("Wallet").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(db *gorm.DB, email string) (*User, error) {
	fmt.Printf("%+v\n", db)
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetMailByIssuer(db *gorm.DB, issuer string) (string, error) {
	var user User

	if err := db.Where("id = ?", issuer).First(&user).Error; err != nil {
		return "", fmt.Errorf("Email not found for issuer %s", issuer)
	}

	UserMail := user.Email
	return UserMail, nil
}

func GetVerficationByMail(db *gorm.DB, email string) (string, error) {
	var verfication Verfication

	if err := db.Where("email = ?", email).First(&verfication).Error; err != nil {
		return "", fmt.Errorf("Verfication code not found for mail %s", email)
	}

	verfication_code := verfication.Verfication
	return verfication_code, nil
}

func GetWalletAddressByIssuer(db *gorm.DB, issuer string) (string, error) {
	var wallet Wallet

	if err := db.Where("user_id = ?", issuer).First(&wallet).Error; err != nil {
		return "", fmt.Errorf("Wallet not found for issuer %s", issuer)
	}

	WalletAddress := wallet.WalletAddress
	return WalletAddress, nil
}

func GetWalletAddress(db *gorm.DB, issuer string) (string, error) {
	var wallet Wallet

	if err := db.Where("user_id = ?", issuer).First(&wallet).Error; err != nil {
		return "", err
	}
	WalletAddress := wallet.WalletAddress
	return WalletAddress, nil
}

func GetWalletbyWalletAddress(db *gorm.DB, WalletAddress string) (*Wallet, error) {
	var wallet Wallet
	if err := db.Where("wallet_address = ?", WalletAddress).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func GetCryptoWallet(db *gorm.DB, wallet_id string) ([]CryptoWallet, error) {
	var cryptoWallets []CryptoWallet

	if err := db.Where("wallet_id = ?", wallet_id).Find(&cryptoWallets).Error; err != nil {
		return nil, err
	}
	return cryptoWallets, nil
}

func GetAllCryptos(db *gorm.DB) ([]Crypto, error) {
	var cryptos []Crypto
	if err := db.Find(&cryptos).Error; err != nil {
		return nil, err
	}
	return cryptos, nil
}
