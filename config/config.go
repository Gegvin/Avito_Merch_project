package config

import (
	"os"
)

type Config struct {
	DBConfig  string
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	// Читаем параметры из переменных окружения
	dbConfig := os.Getenv("DB_CONFIG")
	if dbConfig == "" {
		dbConfig = "host=localhost port=5432 user=postgres password=postgres dbname=avito_merch sslmode=disable"
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "mysecret"
	}
	return &Config{
		DBConfig:  dbConfig,
		JWTSecret: jwtSecret,
	}, nil
}
