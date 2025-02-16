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
	"Avito_Merch_project/internal/middleware"
	"Avito_Merch_project/internal/repository"
	"Avito_Merch_project/internal/services"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func setupSendCoinIntegrationServer() *mux.Router {
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
	coinService := services.NewCoinService(repo)
	handlers.SetAuthService(authService)
	handlers.SetCoinService(coinService)
	handlers.SetRepository(repo)

	r := mux.NewRouter()
	apiPublic := r.PathPrefix("/api").Subrouter()
	apiPublic.HandleFunc("/auth", handlers.AuthHandler).Methods("POST")

	apiPrivate := r.PathPrefix("/api").Subrouter()
	apiPrivate.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	apiPrivate.HandleFunc("/sendCoin", handlers.SendCoinHandler).Methods("POST")
	return r
}

func TestSendCoinIntegration(t *testing.T) {
	router := setupSendCoinIntegrationServer()
	token := getValidToken(t, router)

	payload := map[string]interface{}{
		"toUser": "recipient",
		"amount": 50,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/api/sendCoin", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
