package controllers

import (
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
