package services

import (
	"errors"
	"testing"

	"Avito_Merch_project/internal/models"
)

type fakeRepoForMerch struct {
	purchaseCalled bool
	purchaseError  error
}

func (f *fakeRepoForMerch) CreateUserIfNotExists(username string) error { return nil }
func (f *fakeRepoForMerch) UserExists(username string) (bool, error)    { return false, nil }
func (f *fakeRepoForMerch) CreateUser(username string) error            { return nil }
func (f *fakeRepoForMerch) GetUserInfo(username string) (int, []models.InventoryItem, models.CoinHistory, error) {
	return 1000, []models.InventoryItem{}, models.CoinHistory{}, nil
}
func (f *fakeRepoForMerch) TransferCoins(fromUser, toUser string, amount int) error { return nil }
func (f *fakeRepoForMerch) PurchaseMerch(username, item string, price int) error {
	f.purchaseCalled = true
	if f.purchaseError != nil {
		return f.purchaseError
	}
	if item == "invalid" {
		return errors.New("invalid item")
	}
	return nil
}

func TestBuyMerchValid(t *testing.T) {
	fakeRepo := &fakeRepoForMerch{}
	merchService := NewMerchService(fakeRepo)

	err := merchService.BuyMerch("user1", "t-shirt")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !fakeRepo.purchaseCalled {
		t.Fatal("Expected PurchaseMerch to be called")
	}
}

func TestBuyMerchInvalidItem(t *testing.T) {
	fakeRepo := &fakeRepoForMerch{}
	merchService := NewMerchService(fakeRepo)

	err := merchService.BuyMerch("user1", "invalid")
	if err == nil {
		t.Fatal("Expected error for invalid item")
	}
}

func TestBuyMerch_RepoError(t *testing.T) {
	fakeRepo := &fakeRepoForMerch{purchaseError: errors.New("repo error")}
	merchService := NewMerchService(fakeRepo)

	err := merchService.BuyMerch("user1", "t-shirt")
	if err == nil {
		t.Fatal("Expected error due to repository error")
	}
	if err.Error() != "repo error" {
		t.Fatalf("Expected error 'repo error', got %v", err)
	}
}
