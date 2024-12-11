package service

import (
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/repo"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type BalanceService struct {
	repo *repo.BalanceRepo
}

func NewBalanceService(repo *repo.BalanceRepo) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}

func (b *BalanceService) SetBalance(tx *sqlx.Tx, user *domain.User, value decimal.Decimal) error {
	return b.repo.WithTx(tx).SaveBalance(user.Id, value)
}

func (b *BalanceService) GetBalance(user *domain.User) (*domain.Balance, error) {
	return b.repo.FindBalance(user.Id)
}
