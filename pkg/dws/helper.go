package dws

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateId() int {

	return int(uuid.New().ID())
}

func Convert(amount float64, from, to Currency) (float64, error) {
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
