package usecase

import (
	"fmt"
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/pkg/formater"
	"github.com/SmirnovND/gofermart/internal/pkg/luna"
	"github.com/SmirnovND/gofermart/internal/service"
	"net/http"
)

type OrderUseCase struct {
	orderService      *service.OrderService
	userService       *service.UserService
	processingUseCase *ProcessingUseCase
}

func NewOrderUseCase(
	OrderService *service.OrderService,
	UserService *service.UserService,
	ProcessingUseCase *ProcessingUseCase,
) *OrderUseCase {
	return &OrderUseCase{
		orderService:      OrderService,
		userService:       UserService,
		processingUseCase: ProcessingUseCase,
	}
}

func (o *OrderUseCase) OrdersUpload(w http.ResponseWriter, login string, orderNumber string) {
	validNumber := luna.LunaAlgorithm(orderNumber)
	if !validNumber {
		http.Error(w, "the order number is not valid", http.StatusUnprocessableEntity)
		return
	}

	user, err := o.userService.FindUser(login)
	if err != nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	userId, err := o.orderService.FindUserIdByOrderNumber(orderNumber)
	switch {
	case err != domain.ErrNotFound && userId == user.Id:
		w.WriteHeader(http.StatusOK)
		return
	case err == nil:
		w.WriteHeader(http.StatusConflict)
		return
	case err != domain.ErrNotFound:
		http.Error(w, "FindUserIdByOrderNumber: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = o.orderService.SaveOrder(user.Id, orderNumber)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go func() {
		err = o.processingUseCase.CheckProcessedAndAccrueBalance(orderNumber, user.Id)
		fmt.Println(err)
	}()

	w.WriteHeader(http.StatusAccepted)
	return

}

func (o *OrderUseCase) ListUserOrders(w http.ResponseWriter, login string) {
	user, err := o.userService.FindUser(login)
	if err != nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	orderList, err := o.orderService.ListUserOrders(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(orderList) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	response, err := formater.JSONResponse(orderList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}
