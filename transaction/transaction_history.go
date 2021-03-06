package transaction

// UserTransactionHistory is a struct that represents a transaction
// a negative amount means a debit transaction for single user tx
// a positive amount means a credit transaction for single user tx
// for 2 users tx amount should always be positive
type TransactHistory struct {
	refId  string
	payer  string
	id     string
	amount float64
	payee  string
	rev    bool
	desc   string
}

func (t *TransactHistory) Id() string {
	return t.id
}

func (t *TransactHistory) RefId() string {
	return t.refId
}

func (t *TransactHistory) Desc() string {
	return t.desc
}

func (t *TransactHistory) Amount() float64 {
	return t.amount
}

func (t *TransactHistory) PayeeId() string {
	return t.payee
}

func (t *TransactHistory) PayerId() string {
	return t.payer
}

func (t *TransactHistory) IsTwoPartyTx() bool {
	return len(t.payee) != 0
}

func (t *TransactHistory) IsAccountingTx() bool {
	return t.rev
}

func (t *TransactHistory) IsDebitTx() bool {
	return len(t.payee) == 0 && t.amount > 0
}

func (t *TransactHistory) IsCreditTx() bool {
	return len(t.payee) == 0 && t.amount < 0
}
