package usecase

import (
	"context"
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/pkg/db"
	"github.com/SmirnovND/gofermart/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type ProcessingUseCase struct {
	processingService  *service.ProcessingService
	balanceService     *service.BalanceService
	orderService       *service.OrderService
	transactionManager *db.TransactionManager
}

func NewProcessingUseCase(
	ProcessingService *service.ProcessingService,
	BalanceService *service.BalanceService,
	OrderService *service.OrderService,
	TransactionManager *db.TransactionManager,
) *ProcessingUseCase {
	return &ProcessingUseCase{
		processingService:  ProcessingService,
		balanceService:     BalanceService,
		transactionManager: TransactionManager,
		orderService:       OrderService,
	}
}

func (p *ProcessingUseCase) CheckProcessedAndAccrueBalance(number string, user *domain.User) error {
	order, err := p.processingService.GetOrder(number)
	if err != nil {
		return err
	}

	switch order.Status {
	case service.Processed:
		ctx := context.Background()
		var txErr error
		err = p.transactionManager.Execute(ctx, func(tx *sqlx.Tx) error {
			txErr = p.balanceService.AccrueBalance(tx, user, number, decimal.NewFromFloat(order.Accrual))
			if txErr != nil {
				return txErr
			}

			txErr = p.orderService.ChangeStatus(tx, number, order.Status)
			if txErr != nil {
				return txErr
			}

			return nil
		})
	case service.Invalid, service.Processing:
		ctx := context.Background()
		var txErr error
		err = p.transactionManager.Execute(ctx, func(tx *sqlx.Tx) error {
			txErr = p.orderService.ChangeStatus(tx, number, order.Status)
			if txErr != nil {
				return txErr
			}

			return nil
		})
	}

	return nil
}
