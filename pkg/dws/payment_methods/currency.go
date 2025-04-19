package payment_methods

import (
	"pratikshakuldeep456/digital-wallet-service/pkg/dws"
)

type ICurrency interface {
	Convert(amount float64, fromCurr, toCurr dws.Currency) (float64, error)
}
