package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// отправляет запрос на /api/auth и возвращает полученный токен.
func getValidToken(t *testing.T, router *mux.Router) string {
	authPayload := map[string]string{
		"username": "integrationUser",
		"password": "password",
	}
	body, err := json.Marshal(authPayload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/api/auth", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	if err != nil {
		t.Fatal(err)
	}
	token, ok := resp["token"]
	if !ok || token == "" {
		t.Fatal("Failed to obtain valid token")
	}
	return token
}
