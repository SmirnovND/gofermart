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
	transactionManager *db.TransactionManager
}

func NewProcessingUseCase(
	ProcessingService *service.ProcessingService,
	BalanceService *service.BalanceService,
	TransactionManager *db.TransactionManager,
) *ProcessingUseCase {
	return &ProcessingUseCase{
		processingService:  ProcessingService,
		balanceService:     BalanceService,
		transactionManager: TransactionManager,
	}
}

func (p *ProcessingUseCase) CheckProcessedAndAccrueBalance(number string, user *domain.User) error {
	order, err := p.processingService.GetOrder(number)
	if err != nil {
		return err
	}

	if order.Status == service.Processed {
		ctx := context.Background()
		var txErr error
		err = p.transactionManager.Execute(ctx, func(tx *sqlx.Tx) error {
			txErr = p.balanceService.AccrueBalance(tx, user, number, decimal.NewFromFloat(order.Accrual))
			if txErr != nil {
				return txErr
			}

			return nil
		})

	}

	return nil
}
