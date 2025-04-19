package payment_methods

import "pratikshakuldeep456/digital-wallet-service/pkg/dws"

type Upi struct {
	dws.BasePayment
	UpiId string
}

func NewUpi(upiId string, user *dws.User) *Upi {
	return &Upi{
		BasePayment: dws.BasePayment{
			Id:   dws.GenerateId(),
			User: user,
		},
		UpiId: upiId,
	}
}
func (u *Upi) Pay(amount float64, fromCurr, toCurr dws.Currency) (bool, error) {
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
func (u *Upi) PaymentMethod() dws.PaymentMethod {
	return dws.UpiTransfer
}

func (u *Upi) GetId() int {
	return u.Id
}
