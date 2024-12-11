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

	u.UserUseCase.UserBalance(w, login.(string))
}

func (u *UserController) HandleUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	login := r.Context().Value("login")
	if login == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	u.UserUseCase.UserWithdrawals(w, login.(string))
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

	u.UserUseCase.UserBalanceWithdraw(w, login.(string), withdraw.Number, withdraw.Sum)
}
