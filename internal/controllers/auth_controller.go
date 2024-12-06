package controllers

import (
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/pkg/paramsparser"
	"github.com/SmirnovND/gofermart/internal/usecase"
	"net/http"
)

type AuthController struct {
	AuthUseCase *usecase.AuthUseCase
}

func NewAuthController(AuthUseCase *usecase.AuthUseCase) *AuthController {
	return &AuthController{
		AuthUseCase: AuthUseCase,
	}
}

func (a *AuthController) HandleRegisterJSON(w http.ResponseWriter, r *http.Request) {
	credentials, err := paramsparser.JSONParse[domain.Credentials](w, r)
	if err != nil {
		return
	}
	a.AuthUseCase.Register(w, credentials)
}

func (a *AuthController) HandleLoginJSON(w http.ResponseWriter, r *http.Request) {
	credentials, err := paramsparser.JSONParse[domain.Credentials](w, r)
	if err != nil {
		return
	}
	a.AuthUseCase.Login(w, credentials)
}
