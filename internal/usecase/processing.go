package usecase

import (
	"context"
	"fmt"
	"github.com/SmirnovND/gofermart/internal/pkg/db"
	"github.com/SmirnovND/gofermart/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

type ProcessingUseCase struct {
	processingService  *service.ProcessingService
	balanceService     *service.BalanceService
	orderService       *service.OrderService
	rabbitMqService    *service.RabbitMqService
	transactionManager *db.TransactionManager
}

func NewProcessingUseCase(
	ProcessingService *service.ProcessingService,
	BalanceService *service.BalanceService,
	OrderService *service.OrderService,
	RabbitMqService *service.RabbitMqService,
	TransactionManager *db.TransactionManager,
) *ProcessingUseCase {
	return &ProcessingUseCase{
		processingService:  ProcessingService,
		balanceService:     BalanceService,
		transactionManager: TransactionManager,
		orderService:       OrderService,
		rabbitMqService:    RabbitMqService,
	}
}

func (p *ProcessingUseCase) CheckProcessedAndAccrueBalance(number string, userId int) error {
	order, err := p.processingService.GetOrder(number)
	if err != nil {
		return err
	}

	fmt.Println(order.Status)
	fmt.Println(order.Accrual)
	switch order.Status {
	case service.Processed:
		ctx := context.Background()
		var txErr error
		err = p.transactionManager.Execute(ctx, func(tx *sqlx.Tx) error {
			txErr = p.balanceService.AccrueBalance(tx, userId, number, decimal.NewFromFloat(order.Accrual))
			if txErr != nil {
				return txErr
			}

			txErr = p.orderService.ChangeStatusTx(tx, number, order.Status)
			if txErr != nil {
				return txErr
			}

			txErr = p.orderService.SetAccrualTx(tx, number, order.Accrual)
			if txErr != nil {
				return txErr
			}

			return nil
		})
	case service.Invalid:
		ctx := context.Background()
		var txErr error
		err = p.transactionManager.Execute(ctx, func(tx *sqlx.Tx) error {
			txErr = p.orderService.ChangeStatusTx(tx, number, order.Status)
			if txErr != nil {
				return txErr
			}

			return nil
		})
	default:
		orderModel, err := p.orderService.FindOrderByNumber(number)
		if err != nil {
			return err
		}
		if order.Status != orderModel.Status {
			err = p.orderService.ChangeStatus(number, order.Status)
			if err != nil {
				return err
			}
		}

		fmt.Println("delay")
		return p.rabbitMqService.SendMessageWithDelay(number, userId, 5*time.Second)
	}

	return nil
}
