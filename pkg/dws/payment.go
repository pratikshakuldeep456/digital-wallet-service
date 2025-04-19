package dws

type Payment interface {
	Pay(amount float64, fromCurr, toCurr Currency) (bool, error)
	PaymentMethod() PaymentMethod
	GetId() int
}

type BasePayment struct {
	Id   int
	User *User
}
