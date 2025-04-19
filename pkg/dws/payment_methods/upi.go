package payment_methods

import "pratikshakuldeep456/digital-wallet-service/pkg/dws"

type Upi struct {
	UpiId     string
	Converter dws.CurrencyConverter
}

func NewUpi(upiId string) *Upi {
	return &Upi{
		UpiId:     upiId,
		Converter: &dws.SimpleCurrencyConverter{},
	}
}
func (u *Upi) Pay(amount float64, fromCurr, toCurr dws.Currency) (bool, error) {
	if amount <= 0 {
		return false, nil
	}
	convertedAmount, err := u.Converter.Convert(amount, fromCurr, toCurr)
	if err != nil {
		return false, err
	}
	if convertedAmount == 0 {
		return false, nil
	}
	return true, nil
}
func (u *Upi) PaymentMethod() dws.PaymentMethod {
	return dws.UpiTransfer
}
