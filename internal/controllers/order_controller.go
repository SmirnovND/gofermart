package controllers

import (
	"github.com/SmirnovND/gofermart/internal/pkg/paramsparser"
	"github.com/SmirnovND/gofermart/internal/usecase"
	"net/http"
)

type OrderController struct {
	OrderUseCase *usecase.OrderUseCase
}

func NewOrderController(OrderUseCase *usecase.OrderUseCase) *OrderController {
	return &OrderController{
		OrderUseCase: OrderUseCase,
	}
}

func (o *OrderController) HandleOrdersUpload(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	login := r.Context().Value("login")
	if login == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orderNumber, err := paramsparser.TextParse(w, r)
	if err != nil {
		return
	}

	o.OrderUseCase.OrdersUpload(w, login.(string), orderNumber)
}

func (o *OrderController) HandleListUserOrders(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login")
	if login == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	o.OrderUseCase.ListUserOrders(w, login.(string))
}
