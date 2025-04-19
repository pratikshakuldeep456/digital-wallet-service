package dws

import (
	"errors"
	"fmt"
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
			PaymentMethods: make(map[int][]Payment),
			Transactions:   make(map[int][]*Transaction)}
	})
	return DigitalWalletServiceInstance
}

func (d *DigitalWalletService) CreateUser(name, phone string, profile Profile) *User {
	user := CreateUser(name, phone, profile)
	d.User[user.Id] = user
	return user
}
func (d *DigitalWalletService) UpdateProfile(id int, email, address *string) (*User, error) {
	user := d.User[id]
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.UpdateProfile(email, address)
	return user, nil

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

func (d *DigitalWalletService) FundAccountFromPaymentMethod(userId int, accountId int, paymentMethodId int, amount float64, currency Currency) (bool, error) {
	fmt.Println("Funding account from payment method")
	user := d.User[userId]
	if user == nil {
		return false, errors.New("user not found")
	}
	if user.IsDeleted {
		return false, errors.New("user account is deleted")
	}
	// account := user.Accounts[accountId]
	// if account == nil {
	// 	return false, errors.New("account not found")
	// }
	// fmt.Println("User account:", account)
	paymentMethods, exists := d.PaymentMethods[userId]
	if !exists || len(paymentMethods) == 0 {
		return false, errors.New("no payment methods found for user")
	}

	// Check if the user has any payment methods
	fmt.Println("User payment methods:", d.PaymentMethods)
	// paymentMethods, exists := d.PaymentMethods[userId]
	// if !exists || len(paymentMethods) == 0 {
	// 	return false, errors.New("no payment methods found for user")
	// }
	if amount <= 0 {
		return false, nil
	}
	var selectedPaymentMethod Payment
	for _, method := range paymentMethods {
		if method.GetId() == paymentMethodId {
			selectedPaymentMethod = method
			break
		}
	}
	fmt.Printf("Selected payment method: %v\n", selectedPaymentMethod)
	if selectedPaymentMethod == nil {
		return false, errors.New("payment method not found")
	}

	// Find the account
	var selectedAccount *Account
	for _, account := range user.Accounts {
		if account.Id == accountId {
			selectedAccount = account
			break
		}
	}

	if selectedAccount == nil {
		return false, errors.New("account not found")
	}

	success, err := selectedPaymentMethod.Pay(amount, currency, selectedAccount.Currency)
	if err != nil {
		return false, fmt.Errorf("payment processing failed: %v", err)
	}
	fmt.Printf("Payment successful: %v\n", success)

	if !success {
		return false, errors.New("payment was not successful")
	}

	// Convert currency if needed
	depositAmount := amount
	if currency != selectedAccount.Currency {
		converted, err := Convert(amount, currency, selectedAccount.Currency)
		if err != nil {
			return false, fmt.Errorf("currency conversion failed: %v", err)
		}
		depositAmount = converted
	}

	// Deposit funds to the account
	err = selectedAccount.Deposit(depositAmount)
	if err != nil {
		return false, fmt.Errorf("deposit failed: %v", err)
	}

	// Create and record transaction
	txn := &Transaction{
		Id:              GenerateId(),
		Amount:          depositAmount,
		TransactionType: Credit,
		SrcAccount:      nil, // External source via payment method
		DesAccount:      selectedAccount,
		CreatedAt:       time.Now(),
	}

	// Add transaction to account
	selectedAccount.AddTransaction(txn)

	// Add transaction to user's transaction history
	// Initialize the transactions map if it doesn't exist
	// if d.Transactions == nil {
	// 	d.Transactions = make(map[int][]*Transaction)
	// }

	// Initialize the user's transaction slice if it doesn't exist
	if _, exists := d.Transactions[userId]; !exists {
		d.Transactions[userId] = make([]*Transaction, 0)
	}

	d.Transactions[userId] = append(d.Transactions[userId], txn)

	return true, nil
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
