package service

import (
	"github.com/SmirnovND/gofermart/internal/repo"
	"strconv"
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

// LunaAlgorithm проверяет корректность числа по алгоритму Луна
func (o *OrderService) LunaAlgorithm(orderNumber string) bool {
	var digits []int
	for _, ch := range orderNumber {
		digit, err := strconv.Atoi(string(ch))
		if err != nil {
			return false // если символ не цифра, возвращаем false
		}
		digits = append(digits, digit)
	}

	var sum int
	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i]

		if (len(digits)-i)%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit = digit - 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
