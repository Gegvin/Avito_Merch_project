package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func TestGetUserDataFromToken_Valid(t *testing.T) {
	secret := "testsecret"
	// Создаем JWT-токен с именем "testuser"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Error signing token: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})

	data, ok := getUserDataFromToken(req, secret)
	if !ok {
		t.Fatal("Expected token to be valid")
	}
	if data.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got %q", data.Username)
	}
}

func TestGetUserDataFromToken_Missing(t *testing.T) {
	secret := "testsecret"
	req := httptest.NewRequest("GET", "/", nil)
	_, ok := getUserDataFromToken(req, secret)
	if ok {
		t.Fatal("Expected failure for missing token")
	}
}

func TestGetUserDataFromToken_Invalid(t *testing.T) {
	secret := "testsecret"
	req := httptest.NewRequest("GET", "/", nil)
	// недействителльный токен
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "invalid.token.value",
	})
	_, ok := getUserDataFromToken(req, secret)
	if ok {
		t.Fatal("Expected failure for invalid token")
	}
}
