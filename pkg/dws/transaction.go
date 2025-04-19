package dws

import (
	"time"
)

type Transaction struct {
	Id              int
	Amount          float64
	TransactionType Transactiontype
	SrcAccount      *Account
	DesAccount      *Account
	CreatedAt       time.Time
}

func AddTransaction(txn Transaction) *Transaction {
	return &Transaction{
		Id:              GenerateId(),
		Amount:          txn.Amount,
		TransactionType: txn.TransactionType,
		SrcAccount:      txn.SrcAccount,
		DesAccount:      txn.DesAccount,
		CreatedAt:       time.Now(),
	}
}
