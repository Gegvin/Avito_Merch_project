package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"Avito_Merch_project/internal/services"

	"github.com/stretchr/testify/assert"
)

type fakeAuthService struct{}

func (f *fakeAuthService) Authenticate(username, password string) (string, error) {
	if username == "fail" {
		return "",
			""
	}
	return "dummy-token", nil
}

func (f *fakeAuthService) Register(username, password string) (string, error) {
	if username == "fail" {
		return "",
			""
	}
	return "dummy-token", nil
}

func TestAuthHandler_Success(t *testing.T) {
	//  фиктивный сервис
	services.SetAuthService(&fakeAuthService{})

	payload := map[string]string{
		"username": "testuser",
		"password": "password",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	AuthHandler(rr, req)

	// Ожидаем статус 200
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 for successful auth")
}

func TestAuthHandler_Failure(t *testing.T) {
	//  фиктивный, который возвращает ошибку для имени "fail"
	services.SetAuthService(&fakeAuthService{})
	payload := map[string]string{
		"username": "fail",
		"password": "password",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/auth", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	AuthHandler(rr, req)
	// Ожидаем что статус будет не 200
	assert.NotEqual(t, http.StatusOK, rr.Code, "Expected failure status for auth error")
}

func TestRegisterHandler_Success(t *testing.T) {
	services.SetAuthService(&fakeAuthService{})
	payload := map[string]string{
		"username": "newuser",
		"password": "password",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status 200 for successful registration")
}

func TestRegisterHandler_Failure(t *testing.T) {
	services.SetAuthService(&fakeAuthService{})
	payload := map[string]string{
		"username": "fail",
		"password": "password",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/api/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)
	// Ожидаем, что статус не 200 (ошибка регистрации)
	assert.NotEqual(t, http.StatusOK, rr.Code, "Expected failure status for registration error")
}
