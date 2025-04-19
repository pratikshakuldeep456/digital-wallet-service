package dws

import (
	"sync"
	"time"
)

type DigitalWalletService struct {
	User           map[int]*User
	PaymentMethods map[int][]Payment
	Transactions   map[int][]*Transaction
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

func (d *DigitalWalletService) TransferFunds(sen, rec *Account, amount float64, curr Currency) (error, bool) {

	sourceAmoutn := amount
	if sen.Currency != curr {
		converter, err := Convert(amount, curr, sen.Currency)
		if err != nil {
			return err, false
		}
		sourceAmoutn = converter
	}
	sen.Withdraw(sourceAmoutn)

	recAmount := amount
	if rec.Currency != curr {
		converter, err := Convert(amount, curr, rec.Currency)
		if err != nil {
			return err, false
		}
		recAmount = converter

	}
	rec.Deposit(recAmount)

	// Create a transaction for the sender
	txn := &Transaction{
		Id:              GenerateId(),
		Amount:          sourceAmoutn,
		TransactionType: Debit,
		SrcAccount:      sen,
		DesAccount:      rec,
		CreatedAt:       time.Now(),
	}
	sen.AddTransaction(txn)
	// Create a transaction for the receiver
	txn = &Transaction{
		Id:              GenerateId(),
		Amount:          recAmount,
		TransactionType: Credit,
		SrcAccount:      rec,
		DesAccount:      sen,
		CreatedAt:       time.Now(),
	}
	rec.AddTransaction(txn)
	// Add the transaction to the user's transaction history
	d.Transactions[sen.User.Id] = append(d.Transactions[sen.User.Id], txn)
	d.Transactions[rec.User.Id] = append(d.Transactions[rec.User.Id], txn)
	return nil, true

}

func (d *DigitalWalletService) GetTransactionHistory(userid int) []*Transaction {
	user := d.User[userid]
	if user == nil {
		return nil
	}
	return d.Transactions[userid]
}
