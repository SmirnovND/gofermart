package controllers

import (
	"github.com/SmirnovND/gofermart/internal/pkg/paramsparser"
	"github.com/SmirnovND/gofermart/internal/usecase"
	"net/http"
)

type UserController struct {
	UserUseCase *usecase.UserUseCase
}

func NewUserController(UserUseCase *usecase.UserUseCase) *UserController {
	return &UserController{
		UserUseCase: UserUseCase,
	}
}

func (u *UserController) HandleUserBalance(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login")
	if login == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := u.UserUseCase.UserBalance(login.(string))
	if err != nil {
		http.Error(w, err.Error(), err.Code())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (u *UserController) HandleUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login")
	if login == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := u.UserUseCase.UserWithdrawals(login.(string))
	if err != nil {
		http.Error(w, err.Error(), err.Code())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (u *UserController) HandleUserBalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login")
	if login == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Withdraw struct {
		Number string  `json:"order"`
		Sum    float64 `json:"sum"`
	}

	withdraw, err := paramsparser.JSONParse[Withdraw](w, r)
	if err != nil {
		return
	}

	errDomain := u.UserUseCase.UserBalanceWithdraw(login.(string), withdraw.Number, withdraw.Sum)
	if errDomain != nil {
		http.Error(w, errDomain.Error(), errDomain.Code())
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
