package service

import (
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/repo"
)

type UserService struct {
	repo        *repo.UserRepo
	AuthService *AuthService
}

func NewUserService(repo *repo.UserRepo, AuthService *AuthService) *UserService {
	return &UserService{
		repo:        repo,
		AuthService: AuthService,
	}
}

func (u *UserService) FinUser(login string) (*domain.User, error) {
	return u.FinUser(login)
}

func (u *UserService) SaveUser(login string, pass string) (*domain.User, error) {
	user := &domain.User{}
	user.Login = login

	hash, err := u.AuthService.HashPassword(pass)
	if err != nil {
		return nil, err
	}
	user.PassHash = hash

	err = u.repo.SaveUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
