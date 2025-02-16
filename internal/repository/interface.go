package repository

import "Avito_Merch_project/internal/models"

type RepositoryInterface interface {
	CreateUserIfNotExists(username string) error
	UserExists(username string) (bool, error)
	CreateUser(username string) error
	GetUserInfo(username string) (int, []models.InventoryItem, models.CoinHistory, error)
	TransferCoins(fromUser, toUser string, amount int) error
	PurchaseMerch(username, item string, price int) error
}
