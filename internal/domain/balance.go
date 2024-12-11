package domain

import "github.com/shopspring/decimal"

type Balance struct {
	Value decimal.Decimal `json:"-"`
}
