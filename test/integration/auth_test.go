package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Avito_Merch_project/config"
	"Avito_Merch_project/internal/handlers"
	"Avito_Merch_project/internal/repository"
	"Avito_Merch_project/internal/services"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupAuthIntegrationServer() *mux.Router {
	os.Setenv("DB_CONFIG", "host=localhost port=5433 user=postgres password=postgres dbname=avito_merch sslmode=disable")
	os.Setenv("JWT_SECRET", "mysecret")

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	db, err := repository.NewPostgresDB(cfg.DBConfig)
	if err != nil {
		panic(err)
	}
	repo := repository.NewRepository(db)

	authService := services.NewAuthService(repo, cfg.JWTSecret)
	handlers.SetAuthService(authService)
	handlers.SetRepository(repo)

	r := mux.NewRouter()
	r.HandleFunc("/api/auth", handlers.AuthHandler).Methods("POST")
	r.HandleFunc("/api/register", handlers.RegisterHandler).Methods("POST")
	return r
}

func TestAuthIntegration(t *testing.T) {
	router := setupAuthIntegrationServer()

	authPayload := map[string]string{
		"username": "integrationUser",
		"password": "password",
	}
	body, _ := json.Marshal(authPayload)
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
	assert.NoError(t, err)
	token, ok := resp["token"]
	assert.True(t, ok)
	assert.NotEmpty(t, token)
}
