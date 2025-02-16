package services

import (
	"errors"
	"time"

	"Avito_Merch_project/internal/repository"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	repo      repository.RepositoryInterface
	jwtSecret string
}

func NewAuthService(repo repository.RepositoryInterface, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret}
}

func (s *AuthService) Authenticate(username, password string) (string, error) {
	// Если пользователя нет создаем его автоматически
	if err := s.repo.CreateUserIfNotExists(username); err != nil {
		return "", err
	}
	return s.generateToken(username)
}

func (s *AuthService) Register(username, password string) (string, error) {
	exists, err := s.repo.UserExists(username)
	if err != nil {
		return "", err
	}
	if exists {
		return "", errors.New("пользователь уже существует")
	}
	if err := s.repo.CreateUser(username); err != nil {
		return "", err
	}
	return s.generateToken(username)
}

func (s *AuthService) generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}
