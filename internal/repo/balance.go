package repo

import (
	"fmt"
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type BalanceRepo struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewBalanceRepo(db *sqlx.DB) *BalanceRepo {
	return &BalanceRepo{
		db: db,
	}
}

func (b *BalanceRepo) WithTx(tx *sqlx.Tx) *BalanceRepo {
	return &BalanceRepo{
		db: b.db,
		tx: tx,
	}
}

func (b *BalanceRepo) SaveBalance(userId int, value decimal.Decimal) error {
	exec := b.db.Exec
	if b.tx != nil {
		exec = b.tx.Exec
	}
	query := `INSERT INTO "balance" (user_id, value) VALUES ($1, $2)`
	_, err := exec(query, userId, value)
	if err != nil {
		return fmt.Errorf("error saving balance: %w", err)
	}
	return nil
}

func (b *BalanceRepo) WithdrawBalance(userId int, value decimal.Decimal) error {
	exec := b.db.Exec
	if b.tx != nil {
		exec = b.tx.Exec
	}
	query := `UPDATE "balance" SET value = value - $2, total_points_used = total_points_used + $2 
                 WHERE user_id = $1 AND value >= $2;`
	_, err := exec(query, userId, value)
	if err != nil {
		return fmt.Errorf("error WithdrawBalance: %w", err)
	}
	return nil
}

func (b *BalanceRepo) AccrueBalance(userId int, value decimal.Decimal) error {
	exec := b.db.Exec
	if b.tx != nil {
		exec = b.tx.Exec
	}
	query := `UPDATE "balance" SET value = value + $2 
                 WHERE user_id = $1 AND value >= $2;`
	_, err := exec(query, userId, value)
	if err != nil {
		return fmt.Errorf("error AccrueBalance: %w", err)
	}
	return nil
}

func (b *BalanceRepo) FindBalance(userId int) (*domain.Balance, error) {
	query := `SELECT value, total_points_used FROM "balance" WHERE user_id = $1 LIMIT 1`
	row := b.db.QueryRow(query, userId)

	balance := &domain.Balance{}
	err := row.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, fmt.Errorf("error querying balance: %w", err)
	}

	return balance, nil
}
