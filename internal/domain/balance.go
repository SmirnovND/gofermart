package domain

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type Balance struct {
	Current   decimal.Decimal `json:"-"`
	Withdrawn decimal.Decimal `json:"-"`
}

func (b Balance) MarshalJSON() ([]byte, error) {
	type Alias Balance
	balance := struct {
		Current   interface{} `json:"current"`
		Withdrawn interface{} `json:"withdrawn"`
		Alias
	}{
		Alias: Alias(b),
	}

	if !b.Current.IsZero() {
		if b.Current.Equal(b.Current.Truncate(0)) {
			balance.Current = b.Current.IntPart()
		} else {
			balance.Current, _ = b.Current.Float64()
		}
	} else {
		balance.Current = 0
	}

	if !b.Withdrawn.IsZero() {
		if b.Withdrawn.Equal(b.Withdrawn.Truncate(0)) {
			balance.Withdrawn = b.Withdrawn.IntPart()
		} else {
			balance.Withdrawn, _ = b.Withdrawn.Float64()
		}
	} else {
		balance.Withdrawn = 0
	}

	return json.Marshal(balance)
}
