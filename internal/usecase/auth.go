package usecase

import (
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/service"
	"net/http"
)

type AuthUseCase struct {
	UserService *service.UserService
	AuthService *service.AuthService
}

func NewAuthUseCase(UserService *service.UserService, AuthService *service.AuthService) *AuthUseCase {
	return &AuthUseCase{
		UserService: UserService,
		AuthService: AuthService,
	}
}

func (a *AuthUseCase) Register(w http.ResponseWriter, credentials *domain.Credentials) {
	w.Header().Set("Content-Type", "application/json")

	_, err := a.UserService.FinUser(credentials.Login)
	if err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err != domain.ErrUserNotFound {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	user, err := a.UserService.SaveUser(credentials.Login, credentials.Password)
	if err != nil {
		http.Error(w, "Error save user", http.StatusInternalServerError)
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

func (a *AuthUseCase) Login(w http.ResponseWriter, credentials *domain.Credentials) {
	w.Header().Set("Content-Type", "application/json")

	user, err := a.UserService.FinUser(credentials.Login)
	if err != nil {
		if err == domain.ErrUserNotFound {
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
