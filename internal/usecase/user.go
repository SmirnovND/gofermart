package usecase

import (
	"github.com/SmirnovND/gofermart/internal/pkg/formater"
	"github.com/SmirnovND/gofermart/internal/service"
	"net/http"
)

type UserUseCase struct {
	BalanceService *service.BalanceService
	UserService    *service.UserService
}

func NewUserUseCase(BalanceService *service.BalanceService, UserService *service.UserService) *UserUseCase {
	return &UserUseCase{
		BalanceService: BalanceService,
		UserService:    UserService,
	}
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
