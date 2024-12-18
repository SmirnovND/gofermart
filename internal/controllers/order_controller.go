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
		http.Error(w, "Error TextParse", http.StatusInternalServerError)
		return
	}

	upload, errDomain := o.OrderUseCase.OrdersUpload(login.(string), orderNumber)
	if errDomain != nil {
		http.Error(w, err.Error(), errDomain.Code())
		return
	}

	w.WriteHeader(upload)
}

func (o *OrderController) HandleListUserOrders(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login")
	if login == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := o.OrderUseCase.ListUserOrders(login.(string))
	if err != nil {
		http.Error(w, err.Error(), err.Code())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}
