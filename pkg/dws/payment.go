package dws

type Payment interface {
	Pay(amount float64, fromCurr, toCurr Currency) (bool, error)
	PaymentMethod() PaymentMethod
}

type BasePayment struct {
	Id   int
	User *User
}
