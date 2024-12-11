package repo

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

const (
	debiting = "DEBITING"
	accrual  = "ACCRUAL"
)

type TransactionRepo struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewTransactionRepo(db *sqlx.DB) *TransactionRepo {
	return &TransactionRepo{
		db: db,
	}
}

func (t *TransactionRepo) WithTx(tx *sqlx.Tx) *TransactionRepo {
	return &TransactionRepo{
		db: t.db,
		tx: tx,
	}
}

func (t *TransactionRepo) saveTransaction(userId int, value decimal.Decimal, operation string, number string) error {
	exec := t.db.Exec
	if t.tx != nil {
		exec = t.tx.Exec
	}
	query := `INSERT INTO "transaction" (user_id, value, operation, number) VALUES ($1, $2, $3, $4)`
	_, err := exec(query, userId, value, operation, number)
	if err != nil {
		return fmt.Errorf("error saveTransaction %s: %w", operation, err)
	}
	return nil
}

func (t *TransactionRepo) WithdrawTransaction(userId int, value decimal.Decimal, number string) error {
	return t.saveTransaction(userId, value, debiting, number)
}
