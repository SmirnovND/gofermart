package usecase

import (
	"context"
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/pkg/db"
	"github.com/SmirnovND/gofermart/internal/pkg/formater"
	"github.com/SmirnovND/gofermart/internal/pkg/luna"
	"github.com/SmirnovND/gofermart/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"net/http"
)

type UserUseCase struct {
	balanceService     *service.BalanceService
	userService        *service.UserService
	transactionManager *db.TransactionManager
}

func NewUserUseCase(
	BalanceService *service.BalanceService,
	UserService *service.UserService,
	TransactionManager *db.TransactionManager,
) *UserUseCase {
	return &UserUseCase{
		balanceService:     BalanceService,
		userService:        UserService,
		transactionManager: TransactionManager,
	}
}

func (u *UserUseCase) UserBalanceWithdraw(login string, number string, sum float64) *domain.Error {
	decSum := decimal.NewFromFloat(sum)

	user, err := u.userService.FindUser(login)
	if err != nil {
		return &domain.Error{
			Message:   "user not found",
			CodeValue: http.StatusInternalServerError,
		}
	}

	validNumber := luna.LunaAlgorithm(number)
	if !validNumber {
		return &domain.Error{
			Message:   "the order number is not valid",
			CodeValue: http.StatusUnprocessableEntity,
		}
	}

	ctx := context.Background()
	var txErr error
	err = u.transactionManager.Execute(ctx, func(tx *sqlx.Tx) error {
		txErr = u.balanceService.WithdrawBalance(tx, user, number, decSum)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		if err == domain.ErrInsufficientFunds {
			return &domain.Error{
				Message:   err.Error(),
				CodeValue: http.StatusUnprocessableEntity,
			}
		}
		return &domain.Error{
			Message:   err.Error(),
			CodeValue: http.StatusInternalServerError,
		}
	}

	return nil
}

func (u *UserUseCase) UserWithdrawals(login string) ([]byte, *domain.Error) {
	user, err := u.userService.FindUser(login)
	if err != nil {
		return nil, &domain.Error{
			Message:   "user not found",
			CodeValue: http.StatusInternalServerError,
		}
	}

	withdrawals, err := u.balanceService.GetWithdrawals(user)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, &domain.Error{
				Message:   "",
				CodeValue: http.StatusNoContent,
			}
		}
		return nil, &domain.Error{
			Message:   err.Error(),
			CodeValue: http.StatusInternalServerError,
		}
	}

	response, err := formater.JSONResponse(withdrawals)
	if err != nil {
		return nil, &domain.Error{
			Message:   err.Error(),
			CodeValue: http.StatusInternalServerError,
		}
	}

	return response, nil

}

func (u *UserUseCase) UserBalance(login string) ([]byte, *domain.Error) {
	user, err := u.userService.FindUser(login)
	if err != nil {
		return nil, &domain.Error{
			Message:   "user not found",
			CodeValue: http.StatusInternalServerError,
		}
	}

	balance, err := u.balanceService.GetBalance(user)
	if err != nil {
		return nil, &domain.Error{
			Message:   err.Error(),
			CodeValue: http.StatusInternalServerError,
		}
	}

	response, err := formater.JSONResponse(balance)
	if err != nil {
		return nil, &domain.Error{
			Message:   err.Error(),
			CodeValue: http.StatusInternalServerError,
		}
	}

	return response, nil
}
