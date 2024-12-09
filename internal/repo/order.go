package repo

import (
	"database/sql"
	"fmt"
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/jmoiron/sqlx"
)

type OrderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) FindUserIdByOrderNumber(number string) (int, error) {
	query := `SELECT user_id FROM "order" WHERE number = $1 LIMIT 1`
	row := r.db.QueryRow(query, number)

	var userId int
	err := row.Scan(&userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.ErrNotFound
		}
		return 0, fmt.Errorf("error querying FindUserIdByOrderNumber: %w", err)
	}

	return userId, nil
}

func (r *OrderRepo) SaveOrder(userId int, orderNumber string) error {
	query := `INSERT INTO "order" (user_id, number) VALUES ($1, $2)`
	_, err := r.db.Exec(query, userId, orderNumber)
	if err != nil {
		return fmt.Errorf("error saving order: %w", err)
	}
	return nil
}
