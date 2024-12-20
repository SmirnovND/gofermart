package repo

import (
	"database/sql"
	"fmt"
	"github.com/SmirnovND/gofermart/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type OrderRepo struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewOrderRepo(db *sqlx.DB) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) WithTx(tx *sqlx.Tx) *OrderRepo {
	return &OrderRepo{
		db: r.db,
		tx: tx,
	}
}

func (r *OrderRepo) FindUserOrders(userId int) ([]*domain.Order, error) {
	query := `SELECT number, status, accrual, uploaded_at FROM "order" WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("error querying FindUserOrders: %w", err)
	}
	defer rows.Close()

	var orders []*domain.Order
	for rows.Next() {
		var order domain.Order
		var accrual sql.NullFloat64

		err := rows.Scan(&order.Number, &order.Status, &accrual, &order.UploadedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row in FindUserOrders: %w", err)
		}

		// Преобразование значения accrual
		if accrual.Valid {
			order.Accrual = decimal.NewFromFloat(accrual.Float64)
		} else {
			order.Accrual = decimal.Decimal{} // Пустое значение для Accrual
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows in FindUserOrders: %w", err)
	}

	if len(orders) == 0 {
		return nil, domain.ErrNotFound
	}

	return orders, nil
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
func (r *OrderRepo) FindOrderByNumber(number string) (*domain.Order, error) {
	query := `SELECT number, status, accrual, uploaded_at FROM "order" WHERE number = $1 LIMIT 1`
	row := r.db.QueryRow(query, number)

	order := &domain.Order{}
	var accrual sql.NullFloat64
	err := row.Scan(&order.Number, &order.Status, &accrual, &order.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error querying FindOrderByNumber: %w", err)
	}

	if accrual.Valid {
		order.Accrual = decimal.NewFromFloat(accrual.Float64)
	} else {
		order.Accrual = decimal.Decimal{}
	}

	return order, nil
}

func (r *OrderRepo) SaveOrder(userId int, orderNumber string) error {
	query := `INSERT INTO "order" (user_id, number) VALUES ($1, $2)`
	_, err := r.db.Exec(query, userId, orderNumber)
	if err != nil {
		return fmt.Errorf("error saving order: %w", err)
	}
	return nil
}

func (r *OrderRepo) ChangeStatus(number string, status string) error {
	exec := r.db.Exec
	if r.tx != nil {
		exec = r.tx.Exec
	}
	query := `UPDATE "order" SET status = $2 
                 WHERE number = $1;`
	_, err := exec(query, number, status)
	if err != nil {
		return fmt.Errorf("error ChangeStatusTx: %w", err)
	}
	return nil
}

func (r *OrderRepo) SetAccrual(number string, accrual float64) error {
	exec := r.db.Exec
	if r.tx != nil {
		exec = r.tx.Exec
	}
	query := `UPDATE "order" SET accrual = $2 
                 WHERE number = $1;`
	_, err := exec(query, number, accrual)
	if err != nil {
		return fmt.Errorf("error SetAccrual: %w", err)
	}
	return nil
}
