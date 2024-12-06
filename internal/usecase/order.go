package usecase

import (
	"github.com/SmirnovND/gofermart/internal/service"
	"net/http"
)

type OrderUseCase struct {
	OrderService *service.OrderService
}

func NewOrderUseCase(OrderService *service.OrderService) *OrderUseCase {
	return &OrderUseCase{
		OrderService: OrderService,
	}
}

func (o *OrderUseCase) OrdersLoad(w http.ResponseWriter, login string, orderNumber string) {
	validNumber := o.OrderService.LunaAlgorithm(orderNumber)

	if !validNumber {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
}
