package services

import (
	"errors"
	"testing"

	"Avito_Merch_project/internal/models"

	"github.com/dgrijalva/jwt-go"
)

type fakeRepository struct {
	users      map[string]bool
	forceError bool
}

func (f *fakeRepository) CreateUserIfNotExists(username string) error {
	if f.forceError {
		return errors.New("forced error")
	}
	if f.users == nil {
		f.users = make(map[string]bool)
	}
	if _, exists := f.users[username]; exists {
		return nil
	}
	f.users[username] = true
	return nil
}

func (f *fakeRepository) UserExists(username string) (bool, error) {
	if f.users == nil {
		f.users = make(map[string]bool)
	}
	_, exists := f.users[username]
	return exists, nil
}

func (f *fakeRepository) CreateUser(username string) error {
	if f.forceError {
		return errors.New("forced error")
	}
	if f.users == nil {
		f.users = make(map[string]bool)
	}
	if f.users[username] {
		return errors.New("user exists")
	}
	f.users[username] = true
	return nil
}

func (f *fakeRepository) GetUserInfo(username string) (int, []models.InventoryItem, models.CoinHistory, error) {
	return 1000, []models.InventoryItem{}, models.CoinHistory{}, nil
}

func (f *fakeRepository) TransferCoins(fromUser, toUser string, amount int) error { return nil }
func (f *fakeRepository) PurchaseMerch(username, item string, price int) error    { return nil }

func TestAuthenticate(t *testing.T) {
	repo := &fakeRepository{users: make(map[string]bool)}
	secret := "testsecret"
	authService := NewAuthService(repo, secret)

	token, err := authService.Authenticate("testuser", "password")
	if err != nil {
		t.Fatalf("Authenticate failed: %v", err)
	}
	if token == "" {
		t.Fatal("Expected non-empty token")
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("Error parsing token: %v", err)
	}
	if !parsedToken.Valid {
		t.Fatal("Token is not valid")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Token claims not of type jwt.MapClaims")
	}
	if claims["username"] != "testuser" {
		t.Errorf("Expected username 'testuser', got %v", claims["username"])
	}
}

func TestRegister(t *testing.T) {
	repo := &fakeRepository{users: make(map[string]bool)}
	secret := "testsecret"
	authService := NewAuthService(repo, secret)

	token, err := authService.Register("newuser", "password")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if token == "" {
		t.Fatal("Expected non-empty token")
	}

	_, err = authService.Register("newuser", "password")
	if err == nil {
		t.Fatal("Expected error for existing user")
	}
}

func TestAuthenticate_Error(t *testing.T) {
	repo := &fakeRepository{users: make(map[string]bool), forceError: true}
	secret := "testsecret"
	authService := NewAuthService(repo, secret)

	_, err := authService.Authenticate("testuser", "password")
	if err == nil {
		t.Fatal("Expected error due to forced error in repository")
	}
}
