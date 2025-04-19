package payment_methods

import "pratikshakuldeep456/digital-wallet-service/pkg/dws"

type CreditCard struct {
	CardNumber string
	ExpiryDate string
	CVV        string
	HolderName string
	Converter  dws.CurrencyConverter
}

func NewCreditCard(cardNumber, expiryDate, cvv, holderName string) *CreditCard {
	return &CreditCard{
		CardNumber: cardNumber,
		ExpiryDate: expiryDate,
		CVV:        cvv,
		HolderName: holderName,
		Converter:  &dws.SimpleCurrencyConverter{},
	}
}

func (c *CreditCard) Pay(amount float64, fromCurr, toCurr dws.Currency) (bool, error) {
	if amount <= 0 {
		return false, nil
	}
	convertedAmount, err := c.Converter.Convert(amount, fromCurr, toCurr)
	if err != nil {
		return false, err
	}
	if convertedAmount == 0 {
		return false, nil
	}
	return true, nil
}
func (c *CreditCard) PaymentMethod() dws.PaymentMethod {
	return dws.CreditCard
}
