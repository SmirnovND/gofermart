package domain

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type Withdrawal struct {
	Order       string          `json:"order"`
	Sum         decimal.Decimal `json:"-"`
	ProcessedAt CustomTime      `json:"processed_at"`
}

func (w Withdrawal) MarshalJSON() ([]byte, error) {
	type Alias Withdrawal
	order := struct {
		Sum interface{} `json:"accrual,omitempty"` // Добавляем omitempty
		Alias
	}{
		Alias: Alias(w),
	}

	if !w.Sum.IsZero() {
		if w.Sum.Equal(w.Sum.Truncate(0)) {
			order.Sum = w.Sum.IntPart()
		} else {
			order.Sum, _ = w.Sum.Float64()
		}
	}

	return json.Marshal(order)
}
