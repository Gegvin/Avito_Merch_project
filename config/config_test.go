package config

import "testing"

func TestLoadConfig(t *testing.T) {
	// переменные окружения
	t.Setenv("DB_CONFIG", "host=localhost port=5433 user=postgres password=postgres dbname=avito_merch sslmode=disable")
	t.Setenv("JWT_SECRET", "test_secret")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	expectedDB := "host=localhost port=5433 user=postgres password=postgres dbname=avito_merch sslmode=disable"
	if cfg.DBConfig != expectedDB {
		t.Errorf("Expected DBConfig to be %q, got %q", expectedDB, cfg.DBConfig)
	}
	if cfg.JWTSecret != "test_secret" {
		t.Errorf("Expected JWTSecret to be 'test_secret', got %q", cfg.JWTSecret)
	}
}
