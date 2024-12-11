package repo

import (
	"database/sql"
	"fmt"
	"github.com/SmirnovND/gofermart/internal/domain"
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

func (t *TransactionRepo) GetWithdrawals(userId int) ([]*domain.Withdrawal, error) {
	query := `SELECT value, number, updated_at 
				FROM "transaction" 
			   WHERE user_id = $1 AND operation = $2 
			ORDER BY updated_at DESC`
	rows, err := t.db.Query(query, userId, debiting)
	if err != nil {
		return nil, fmt.Errorf("error querying GetWithdrawals: %w", err)
	}
	defer rows.Close()

	var withdrawals []*domain.Withdrawal
	for rows.Next() {
		var withdrawal domain.Withdrawal
		var value sql.NullFloat64

		err := rows.Scan(&value, &withdrawal.Order, &withdrawal.ProcessedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row in GetWithdrawals: %w", err)
		}

		if value.Valid {
			withdrawal.Sum = decimal.NewFromFloat(value.Float64)
		} else {
			withdrawal.Sum = decimal.Decimal{}
		}

		withdrawals = append(withdrawals, &withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows in GetWithdrawals: %w", err)
	}

	if len(withdrawals) == 0 {
		return nil, domain.ErrNotFound
	}

	return withdrawals, nil
}
