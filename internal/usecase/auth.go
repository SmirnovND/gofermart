package usecase

import (
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/repo"
)

type AuthUseCase struct {
	UserRepo *repo.UserRepo
}

func NewAuth(UserRepo *repo.UserRepo) *AuthUseCase {
	return &AuthUseCase{
		UserRepo: UserRepo,
	}
}

func (a *AuthUseCase) Register(credentials *domain.Credentials) {

}
