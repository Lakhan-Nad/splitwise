package core

// negative indicates user owes
// positive indicates user is owed
type BalanceAccount struct {
	AccountId string
	Amount    float64
}

type Payable struct {
	PayerId string
	PayeeId string
	Amount  float64
}
