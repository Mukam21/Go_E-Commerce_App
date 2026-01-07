package repository

import (
	"github.com/Mukam21/Go_E-Commerce_App/internal/domain"
	"github.com/Mukam21/Go_E-Commerce_App/internal/dto"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreatePayment(payment *domain.Payment) error
	FindOrders(uId uint) ([]domain.OrderItem, error)
	FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error)
}

type transactionStorage struct {
	db *gorm.DB
}

func (t transactionStorage) CreatePayment(payment *domain.Payment) error {
	panic("sd")
}

func (t transactionStorage) FindOrders(uId uint) ([]domain.OrderItem, error) {
	panic("sd")
}

func (t transactionStorage) FindOrderById(uId uint, id uint) (dto.SellerOrderDetails, error) {
	panic("sd")
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionStorage{
		db: db,
	}
}
