package dws

import (
	"errors"
	"sync"
	"time"
)

type DigitalWalletService struct {
	User           map[int]*User
	PaymentMethods map[int][]Payment
}

var (
	DigitalWalletServiceInstance *DigitalWalletService
	Once                         sync.Once
)

func GetDigitalWalletService() *DigitalWalletService {

	Once.Do(func() {
		DigitalWalletServiceInstance = &DigitalWalletService{
			User:           make(map[int]*User),
			PaymentMethods: make(map[int][]Payment)}
	})
	return DigitalWalletServiceInstance
}

func (d *DigitalWalletService) CreateUser(name, phone string, profile Profile) *User {
	user := CreateUser(name, phone, profile)
	d.User[user.Id] = user
	return user
}
func (d *DigitalWalletService) UpdateProfile(id int, email, address *string) *User {
	user := d.User[id]
	if user == nil {
		return nil
	}
	user.UpdateProfile(email, address)
	return user

}

func (d *DigitalWalletService) AddPaymentMethod(userid int, payment Payment) {
	user := d.User[userid]
	if user == nil {
		return
	}

	d.PaymentMethods[userid] = append(d.PaymentMethods[userid], payment)
}

func (d *DigitalWalletService) RemovePaymentMethod(userid int, payment PaymentMethod) {
	user := d.User[userid]
	if user == nil {
		return
	}
	payments := d.PaymentMethods[userid]
	for i, p := range payments {
		if p.PaymentMethod() == payment {
			d.PaymentMethods[userid] = append(payments[:i], payments[i+1:]...)
			break
		}
	}
}

func (d *DigitalWalletService) GetPaymentMethods(userid int) []Payment {
	user := d.User[userid]
	if user == nil {
		return nil
	}
	return d.PaymentMethods[userid]
}

func (d *DigitalWalletService) TransferFunds(senderID int,
	receiverID int, amount float64, payment Payment, currfrom Currency, currto Currency) (error, bool) {
	sender, receiver := d.User[senderID], d.User[receiverID]
	if sender == nil || receiver == nil {
		return errors.New("invalid user(s)"), false
	}
	data, err := payment.Pay(amount, currfrom, currto)
	if err != nil || !data {
		return errors.New("payment failed"), false

	}

	//add transactions
	txn := Transaction{
		Id:              GenerateId(),
		Amount:          amount,
		TransactionType: Debit,
		SrcAccount:      sender.Accounts[0], // Simplified
		DesAccount:      receiver.Accounts[0],
		CreatedAt:       time.Now(),
	}

	convertedAmount, err1 := payment.CurrencyConverter().Convert(amount, currfrom, currto)
	if err1 != nil {
		return err1, false
	}

	sender.Accounts[0].Balance -= amount
	receiver.Accounts[0].Balance += convertedAmount
	sender.Accounts[0].Transactions = append(sender.Accounts[0].Transactions, &txn)
	receiver.Accounts[0].Transactions = append(receiver.Accounts[0].Transactions, &txn)

	return nil, true

}
