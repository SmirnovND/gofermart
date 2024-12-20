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

func (o *OrderUseCase) OrdersUpload(login string, orderNumber string) (int, *domain.Error) {
	validNumber := luna.LunaAlgorithm(orderNumber)
	if !validNumber {
		return 0, &domain.Error{
			Message:   "the order number is not valid",
			CodeValue: http.StatusUnprocessableEntity,
		}
	}

	user, err := o.userService.FindUser(login)
	if err != nil {
		return 0, &domain.Error{
			Message:   "user not found",
			CodeValue: http.StatusInternalServerError,
		}
	}

	userId, err := o.orderService.FindUserIdByOrderNumber(orderNumber)
	switch {
	case err != domain.ErrNotFound && userId == user.Id:
		return http.StatusOK, nil
	case err == nil:
		return 0, &domain.Error{
			Message:   "the order is registered to another user",
			CodeValue: http.StatusConflict,
		}
	case err != domain.ErrNotFound:
		return 0, &domain.Error{
			Message:   "FindUserIdByOrderNumber: " + err.Error(),
			CodeValue: http.StatusInternalServerError,
		}
	}

	err = o.orderService.SaveOrder(user.Id, orderNumber)
	if err != nil {
		return 0, &domain.Error{
			Message:   err.Error(),
			CodeValue: http.StatusInternalServerError,
		}
	}

	go func() {
		err = o.processingUseCase.CheckProcessedAndAccrueBalance(orderNumber, user.Id)
		fmt.Println(err)
	}()

	return http.StatusAccepted, nil

}

func (o *OrderUseCase) ListUserOrders(login string) ([]byte, *domain.Error) {
	user, err := o.userService.FindUser(login)
	if err != nil {
		return nil, &domain.Error{
			Message:   "user not found",
			CodeValue: http.StatusInternalServerError,
		}
	}

	orderList, err := o.orderService.ListUserOrders(user.Id)
	if err != nil {
		return nil, &domain.Error{
			Message:   err.Error(),
			CodeValue: http.StatusInternalServerError,
		}
	}

	if len(orderList) == 0 {
		return nil, &domain.Error{
			Message:   "",
			CodeValue: http.StatusNoContent,
		}
	}

	response, err := formater.JSONResponse(orderList)
	if err != nil {
		return nil, &domain.Error{
			Message:   err.Error(),
			CodeValue: http.StatusInternalServerError,
		}
	}

	return response, nil
}
