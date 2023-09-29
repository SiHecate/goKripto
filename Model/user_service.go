package model

import (
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
	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetWalletAddress(db *gorm.DB, issuer string) (string, error) {
	var wallet Wallet

	if err := db.Where("user_id = ?", issuer).First(&wallet).Error; err != nil {
		return "", err
	}

	WalletAddress := wallet.WalletAddress
	return WalletAddress, nil
}

func GetCryptoWalletsByWalletAddress(db *gorm.DB, walletAddress string) ([]CryptoWallet, error) {
	var cryptoWallets []CryptoWallet

	if err := db.Where("walletAddress = ?", walletAddress).Find(&cryptoWallets).Error; err != nil {
		return nil, err
	}
	return cryptoWallets, nil
}
