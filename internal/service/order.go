package service

import (
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/repo"
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

func (o *OrderService) ListUserOrders(userId int) ([]*domain.Order, error) {
	return o.orderRepo.FindUserOrders(userId)
}
