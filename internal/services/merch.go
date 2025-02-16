package services

import (
	"Avito_Merch_project/internal/repository"
	"errors"
)

var MerchCatalog = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}

type MerchService struct {
	repo repository.RepositoryInterface
}

func NewMerchService(repo repository.RepositoryInterface) *MerchService {
	return &MerchService{repo: repo}
}

func (s *MerchService) BuyMerch(username, item string) error {
	price, ok := MerchCatalog[item]
	if !ok {
		return errors.New("недопустимый товар")
	}
	return s.repo.PurchaseMerch(username, item, price)
}
