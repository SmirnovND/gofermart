package service

import (
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/repo"
	"github.com/jmoiron/sqlx"
)

type UserService struct {
	repo        *repo.UserRepo
	authService *AuthService
}

func NewUserService(repo *repo.UserRepo, AuthService *AuthService) *UserService {
	return &UserService{
		repo:        repo,
		authService: AuthService,
	}
}

func (u *UserService) FindUser(login string) (*domain.User, error) {
	return u.repo.FindUser(login)
}

func (u *UserService) SaveUser(tx *sqlx.Tx, login string, pass string) (*domain.User, error) {
	user := &domain.User{}
	user.Login = login

	hash, err := u.authService.HashPassword(pass)
	if err != nil {
		return nil, err
	}
	user.PassHash = hash

	err = u.repo.WithTx(tx).SaveUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
