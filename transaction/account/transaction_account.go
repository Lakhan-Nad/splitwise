package account

type TransactionAccount struct {
	userId      string
	balance     float64
	inUse       bool
	lastUpdated int64
}

func (t *TransactionAccount) UserId() string {
	return t.userId
}

func (t *TransactionAccount) Balance() float64 {
	return t.balance
}

func (t *TransactionAccount) InUse() bool {
	return t.inUse
}

func (t *TransactionAccount) LastUpdated() int64 {
	return t.lastUpdated
}
