package payment_methods

import (
	"pratikshakuldeep456/digital-wallet-service/pkg/dws"
	// "pratikshakuldeep456/digital-wallet-service/pkg/dws/payment_methods"
)

type Bank struct {
	BankName      string
	HolderName    string
	AccountNumber string
	IFSCCode      string
	converter     *dws.SimpleCurrencyConverter
}

func NewBank(BankName string,
	HolderName string,
	AccountNumber string,
	IFSCCode string) *Bank {
	return &Bank{
		BankName:      BankName,
		HolderName:    HolderName,
		AccountNumber: AccountNumber,
		IFSCCode:      IFSCCode,
		converter:     &dws.SimpleCurrencyConverter{},
	}

}
func (b *Bank) Pay(amount float64, fromCurr, toCurr dws.Currency) (bool, error) {
	if amount <= 0 {
		return false, nil
	}
	convertedAmount, err := b.converter.Convert(amount, fromCurr, toCurr)
	if err != nil {
		return false, err
	}

	if convertedAmount == 0 {
		return false, nil
	}

	return true, nil
}

func (b *Bank) PaymentMethod() dws.PaymentMethod {
	return dws.BankTransfer
}
