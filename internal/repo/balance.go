package repo

import (
	"fmt"
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
	fmt.Println(userId)
	fmt.Println(value)
	query := `INSERT INTO "balance" (user_id, value) VALUES ($1, $2)`
	_, err := exec(query, userId, value)
	if err != nil {
		return fmt.Errorf("error saving balance: %w", err)
	}
	return nil
}
