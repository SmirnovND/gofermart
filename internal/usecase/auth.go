package usecase

import (
	"context"
	"fmt"
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/pkg/db"
	"github.com/SmirnovND/gofermart/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"net/http"
)

type AuthUseCase struct {
	UserService        *service.UserService
	BalanceService     *service.BalanceService
	AuthService        *service.AuthService
	TransactionManager *db.TransactionManager
}

func NewAuthUseCase(
	UserService *service.UserService,
	AuthService *service.AuthService,
	BalanceService *service.BalanceService,
	TransactionManager *db.TransactionManager,
) *AuthUseCase {
	return &AuthUseCase{
		UserService:        UserService,
		BalanceService:     BalanceService,
		AuthService:        AuthService,
		TransactionManager: TransactionManager,
	}
}

func (a *AuthUseCase) Register(w http.ResponseWriter, credentials *domain.Credentials) {
	w.Header().Set("Content-Type", "application/json")

	_, err := a.UserService.FindUser(credentials.Login)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err != domain.ErrNotFound {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	var user *domain.User
	var txErr error

	err = a.TransactionManager.Execute(ctx, func(tx *sqlx.Tx) error {
		user, txErr = a.UserService.SaveUser(tx, credentials.Login, credentials.Password)
		if txErr != nil {
			return txErr
		}

		txErr = a.BalanceService.SetBalance(tx, user, decimal.NewFromInt(0))
		if txErr != nil {
			return txErr
		}

		return nil
	})

	if err != nil {
		// Обработка ошибки сохранения пользователя
		http.Error(w, fmt.Sprintf("Error saving user: %v", err), http.StatusInternalServerError)
		return
	}

	token, err := a.AuthService.GenerateToken(user.Login)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating token: %v", err), http.StatusInternalServerError)
		return
	}

	a.AuthService.SetResponseAuthData(w, token)

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}

func (a *AuthUseCase) Login(w http.ResponseWriter, credentials *domain.Credentials) {
	w.Header().Set("Content-Type", "application/json")

	user, err := a.UserService.FindUser(credentials.Login)
	if err != nil {
		if err == domain.ErrNotFound {
			http.Error(w, "Error", http.StatusUnauthorized)
			return
		} else {
			http.Error(w, "Error", http.StatusInternalServerError)
			return
		}
	}

	// Генерируем токен
	passValid := a.AuthService.CheckPasswordHash(credentials.Password, user.PassHash)
	if !passValid {
		http.Error(w, "Error", http.StatusUnauthorized)
		return
	}

	// Генерируем токен
	token, err := a.AuthService.GenerateToken(user.Login)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	a.AuthService.SetResponseAuthData(w, token)

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}
