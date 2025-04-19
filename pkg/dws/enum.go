package dws

type Transactiontype string

const (
	Credit Transactiontype = "credit"
	Debit  Transactiontype = "debit"
)

type PaymentMethod string

const (
	CreditCard   PaymentMethod = "creditcard"
	BankTransfer PaymentMethod = "banktransfer"
	UpiTransfer  PaymentMethod = "upitransfer"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	INR Currency = "INR"
)

var ExchangeRate = map[string]float64{
	"USD->INR": 83.0,
	"INR->USD": 0.012,
	"USD->EUR": 0.92,
	"EUR->USD": 1.09,
	"EUR->INR": 90.5,
	"INR->EUR": 0.011,
	"GBP->INR": 104.0,
	"INR->GBP": 0.0096,
	// ... add more as needed
}
