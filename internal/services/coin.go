package services

import (
	"Avito_Merch_project/internal/repository"
)

type CoinService struct {
	repo repository.RepositoryInterface
}

func NewCoinService(repo repository.RepositoryInterface) *CoinService {
	return &CoinService{repo: repo}
}

func (s *CoinService) SendCoins(fromUser, toUser string, amount int) error {
	if amount <= 0 {
		return repository.ErrInvalidAmount
	}
	return s.repo.TransferCoins(fromUser, toUser, amount)
}
