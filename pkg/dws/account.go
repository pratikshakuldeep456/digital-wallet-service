package dws

import "sync"

type Account struct {
	Id int
	// User         *User //not needed
	AccountNo    string
	Balance      float64
	Transactions []*Transaction
	Mu           sync.Mutex
}
