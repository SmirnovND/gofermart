package service

import (
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/SmirnovND/gofermart/internal/repo"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type BalanceService struct {
	repo            *repo.BalanceRepo
	transactionRepo *repo.TransactionRepo
}

func NewBalanceService(repo *repo.BalanceRepo, transactionRepo *repo.TransactionRepo) *BalanceService {
	return &BalanceService{
		repo:            repo,
		transactionRepo: transactionRepo,
	}
}

func (b *BalanceService) SetBalance(tx *sqlx.Tx, user *domain.User, value decimal.Decimal) error {
	return b.repo.WithTx(tx).SaveBalance(user.Id, value)
}

func (b *BalanceService) GetBalance(user *domain.User) (*domain.Balance, error) {
	return b.repo.FindBalance(user.Id)
}

func (b *BalanceService) WithdrawBalance(tx *sqlx.Tx, user *domain.User, number string, decSum decimal.Decimal) error {
	balance, err := b.GetBalance(user)
	if err != nil {
		return err
	}

	if balance.Current.Cmp(decSum) < 0 {
		return domain.ErrInsufficientFunds
	}

	err = b.transactionRepo.WithTx(tx).WithdrawTransaction(user.Id, decSum, number)
	if err != nil {
		return err
	}

	err = b.repo.WithTx(tx).WithdrawBalance(user.Id, decSum)
	if err != nil {
		return err
	}

	return nil
}

func (b *BalanceService) GetWithdrawals(user *domain.User) ([]*domain.Withdrawal, error) {
	return b.transactionRepo.GetWithdrawals(user.Id)
}

func (b *BalanceService) AccrueBalance(tx *sqlx.Tx, userId int, number string, decSum decimal.Decimal) error {
	err := b.transactionRepo.WithTx(tx).AccrueTransaction(userId, decSum, number)
	if err != nil {
		return err
	}

	err = b.repo.WithTx(tx).AccrueBalance(userId, decSum)
	if err != nil {
		return err
	}

	return nil
}
