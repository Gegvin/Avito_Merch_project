package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJWTMiddleware_NoToken(t *testing.T) {
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	// Создаем middleware с тестовым секретом
	middleware := JWTMiddleware("testsecret")
	handler := middleware(next)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if nextCalled {
		t.Fatal("Expected middleware to block request without token")
	}
	// Можно проверить статус ошибки (например, 401 Unauthorized)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected status 401, got %d", rr.Code)
	}
}
