package domain

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type Order struct {
	Number     string          `json:"number"`
	Status     string          `json:"status"`
	Accrual    decimal.Decimal `json:"-"`
	UploadedAt CustomTime      `json:"uploaded_at"`
}

type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Format(time.RFC3339) + `"`), nil
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

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("cannot scan type %T into CustomTime", value)
	}
	ct.Time = t
	return nil
}
