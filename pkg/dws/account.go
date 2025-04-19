package dws

import "sync"

type Account struct {
	Id           int
	User         *User //not needed
	AccountNo    string
	Balance      float64
	Currency     Currency
	Transactions []*Transaction
	Mu           sync.Mutex
}

func (a *Account) Deposit(amount float64) error {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	if amount <= 0 {
		return nil
	}

	a.Balance += amount
	return nil
}
func (a *Account) Withdraw(amount float64) error {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	if amount <= 0 {
		return nil
	}

	if a.Balance < amount {
		return nil
	}

	a.Balance -= amount
	return nil
}

func (a *Account) Transfer(amount float64, to *Account) error {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	if amount <= 0 {
		return nil
	}

	if a.Balance < amount {
		return nil
	}

	a.Balance -= amount
	to.Balance += amount

	return nil
}
func (a *Account) AddTransaction(t *Transaction) {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	a.Transactions = append(a.Transactions, t)
}

func (a *Account) GetTransactions() []*Transaction {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	return a.Transactions
}

func (a *Account) GetBalance() float64 {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	return a.Balance
}
