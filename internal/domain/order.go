package domain

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type Order struct {
	Number     string          `json:"number"`
	Status     string          `json:"status"`
	Accrual    decimal.Decimal `json:"-"`
	UploadedAt CustomTime      `json:"uploaded_at"`
}

func (o Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	order := struct {
		Accrual interface{} `json:"accrual,omitempty"` // Добавляем omitempty
		Alias
	}{
		Alias: Alias(o),
	}

	// Если Accrual — это ноль, не добавляем его в JSON
	if !o.Accrual.IsZero() {
		if o.Accrual.Equal(o.Accrual.Truncate(0)) {
			// Если Accrual — целое число, сохраняем его как int
			order.Accrual = o.Accrual.IntPart()
		} else {
			order.Accrual, _ = o.Accrual.Float64()
		}
	}

	return json.Marshal(order)
}
