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
	BalanceService     *service.BalanceService
	UserService        *service.UserService
	OrderService       *service.OrderService
	TransactionManager *db.TransactionManager
}

func NewUserUseCase(
	BalanceService *service.BalanceService,
	UserService *service.UserService,
	OrderService *service.OrderService,
	TransactionManager *db.TransactionManager,
) *UserUseCase {
	return &UserUseCase{
		BalanceService:     BalanceService,
		UserService:        UserService,
		OrderService:       OrderService,
		TransactionManager: TransactionManager,
	}
}

func (u *UserUseCase) UserBalanceWithdraw(w http.ResponseWriter, login string, number string, sum float64) {
	decSum := decimal.NewFromFloat(sum)

	user, err := u.UserService.FindUser(login)
	if err != nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	validNumber := luna.LunaAlgorithm(number)
	if !validNumber {
		http.Error(w, "the order number is not valid", http.StatusUnprocessableEntity)
		return
	}

	ctx := context.Background()
	var txErr error
	err = u.TransactionManager.Execute(ctx, func(tx *sqlx.Tx) error {
		txErr = u.BalanceService.BalanceWithdraw(tx, user, number, decSum)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		if err == domain.ErrInsufficientFunds {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func (u *UserUseCase) UserBalance(w http.ResponseWriter, login string) {
	user, err := u.UserService.FindUser(login)
	if err != nil {
		http.Error(w, "user not found", http.StatusInternalServerError)
		return
	}

	balance, err := u.BalanceService.GetBalance(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := formater.JSONResponse(balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}
