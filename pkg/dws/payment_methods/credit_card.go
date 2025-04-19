package payment_methods

import "pratikshakuldeep456/digital-wallet-service/pkg/dws"

type CreditCard struct {
	dws.BasePayment
	CardNumber string
	ExpiryDate string
	CVV        string
	HolderName string
}

func NewCreditCard(cardNumber, expiryDate, cvv, holderName string, user *dws.User) *CreditCard {
	return &CreditCard{
		BasePayment: dws.BasePayment{
			Id:   dws.GenerateId(),
			User: user,
		},
		CardNumber: cardNumber,
		ExpiryDate: expiryDate,
		CVV:        cvv,
		HolderName: holderName,
	}
}

func (c *CreditCard) Pay(amount float64, fromCurr, toCurr dws.Currency) (bool, error) {
	if amount <= 0 {
		return false, nil
	}
	convertedAmount, err := dws.Convert(amount, fromCurr, toCurr)
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

func (c *CreditCard) GetId() int {
	return c.Id
}
