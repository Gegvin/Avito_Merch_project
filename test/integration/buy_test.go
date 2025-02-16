package integration

import (
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

func setupBuyIntegrationServer() *mux.Router {
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

	r := mux.NewRouter()
	apiPublic := r.PathPrefix("/api").Subrouter()
	apiPublic.HandleFunc("/auth", handlers.AuthHandler).Methods("POST")
	apiPublic.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")

	apiPrivate := r.PathPrefix("/api").Subrouter()
	apiPrivate.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	apiPrivate.HandleFunc("/buy/{item}", handlers.BuyHandler).Methods("GET")

	// Устанавливаем сервисы
	authService := services.NewAuthService(repo, cfg.JWTSecret)
	handlers.SetAuthService(authService)
	handlers.SetRepository(repo)
	merchService := services.NewMerchService(repo)
	handlers.SetMerchService(merchService)

	return r
}

func TestBuyIntegration(t *testing.T) {
	router := setupBuyIntegrationServer()
	token := getValidToken(t, router)

	req, err := http.NewRequest("GET", "/api/buy/t-shirt", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Передаем токен в заголовке и cookie
	req.Header.Set("Authorization", "Bearer "+token)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
