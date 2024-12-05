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
	parseCredentials, err := paramsparser.JSONParse[domain.Credentials](w, r)
	if err != nil {
		return
	}
	a.AuthUseCase.Register(parseCredentials)
	//JSONResponse, err := serverSaver.FindAndResponseAsJSON(parseMetric, mc.ServiceCollector, w)
	//if err != nil {
	//	return
	//}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("hello"))
}
