package dws

import (
	"fmt"
)

type Payment interface {
	Pay(amount float64, fromCurr, toCurr Currency) (bool, error)
	PaymentMethod() PaymentMethod
	CurrencyConverter() CurrencyConverter
}
type CurrencyConverter interface {
	Convert(amount float64, from, to Currency) (float64, error)
}
type SimpleCurrencyConverter struct{}

func (s *SimpleCurrencyConverter) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	key := string(from) + "->" + string(to)
	rate, ok := ExchangeRate[key]
	if !ok {
		return 0, fmt.Errorf("no conversion rate from %s to %s", from, to)
	}

	return amount * rate, nil
}
