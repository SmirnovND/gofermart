package service

import (
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/repo"
	"github.com/jmoiron/sqlx"
)

type OrderService struct {
	orderRepo *repo.OrderRepo
}

func NewOrderService(OrderRepo *repo.OrderRepo) *OrderService {
	return &OrderService{
		orderRepo: OrderRepo,
	}
}

func (o *OrderService) SaveOrder(userId int, OrderNumber string) error {
	return o.orderRepo.SaveOrder(userId, OrderNumber)
}

func (o *OrderService) FindUserIdByOrderNumber(orderNumber string) (int, error) {
	return o.orderRepo.FindUserIdByOrderNumber(orderNumber)
}

func (o *OrderService) FindOrderByNumber(orderNumber string) (*domain.Order, error) {
	return o.orderRepo.FindOrderByNumber(orderNumber)
}

func (o *OrderService) ListUserOrders(userId int) ([]*domain.Order, error) {
	return o.orderRepo.FindUserOrders(userId)
}

func (o *OrderService) ChangeStatusTx(tx *sqlx.Tx, number string, status string) error {
	return o.orderRepo.WithTx(tx).ChangeStatus(number, status)
}

func (o *OrderService) SetAccrualTx(tx *sqlx.Tx, number string, accrual float64) error {
	return o.orderRepo.WithTx(tx).SetAccrual(number, accrual)
}

func (o *OrderService) ChangeStatus(number string, status string) error {
	return o.orderRepo.ChangeStatus(number, status)
}
