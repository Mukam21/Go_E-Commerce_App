package service

import (
	"github.com/Mukam21/Go_E-Commerce_App/internal/domain"
	"github.com/Mukam21/Go_E-Commerce_App/internal/dto"
	"github.com/Mukam21/Go_E-Commerce_App/internal/helper"
	"github.com/Mukam21/Go_E-Commerce_App/internal/repository"
)

type TransactionService struct {
	Repo repository.TransactionRepository
	Auth helper.Auth
}

func (s TransactionService) GetOrders(u domain.User) ([]domain.OrderItem, error) {
	orders, err := s.Repo.FindOrders(u.ID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s TransactionService) GetOrderDetails(u domain.User, id uint) (dto.SellerOrderDetails, error) {

	order, err := s.Repo.FindOrderById(u.ID, id)
	if err != nil {
		return dto.SellerOrderDetails{}, err
	}

	return order, nil
}

func NewTransactionService(r repository.TransactionRepository, auth helper.Auth) *TransactionService {
	return &TransactionService{
		Repo: r,
		Auth: auth,
	}
}
