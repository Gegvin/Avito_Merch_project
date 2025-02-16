package services

import (
	"errors"
	"testing"

	"Avito_Merch_project/internal/models"
)

type fakeRepoForCoin struct {
	transferCalled bool
	forceError     bool
}

func (f *fakeRepoForCoin) CreateUserIfNotExists(username string) error { return nil }
func (f *fakeRepoForCoin) UserExists(username string) (bool, error)    { return false, nil }
func (f *fakeRepoForCoin) CreateUser(username string) error            { return nil }
func (f *fakeRepoForCoin) GetUserInfo(username string) (int, []models.InventoryItem, models.CoinHistory, error) {
	return 1000, []models.InventoryItem{}, models.CoinHistory{}, nil
}
func (f *fakeRepoForCoin) TransferCoins(fromUser, toUser string, amount int) error {
	if f.forceError {
		return errors.New("forced error")
	}
	f.transferCalled = true
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	return nil
}
func (f *fakeRepoForCoin) PurchaseMerch(username, item string, price int) error { return nil }

func TestSendCoins(t *testing.T) {
	fakeRepo := &fakeRepoForCoin{}
	coinService := NewCoinService(fakeRepo)

	// Проверка передачи 0 монет.
	err := coinService.SendCoins("user1", "user2", 0)
	if err == nil {
		t.Fatal("Expected error for zero amount")
	}

	// Корректная передача монет.
	err = coinService.SendCoins("user1", "user2", 50)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !fakeRepo.transferCalled {
		t.Fatal("Expected TransferCoins to be called")
	}
}

func TestSendCoins_RepoError(t *testing.T) {
	fakeRepo := &fakeRepoForCoin{forceError: true}
	coinService := NewCoinService(fakeRepo)

	err := coinService.SendCoins("user1", "user2", 50)
	if err == nil {
		t.Fatal("Expected error due to forced error in repository")
	}
}
