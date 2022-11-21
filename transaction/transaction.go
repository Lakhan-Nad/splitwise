package transaction

import (
	"fmt"

	Account "workat.tech/splitwise/transaction/account"
)

var (
	idDivider      = "@TxId"
	revIdDivider   = "@Rev"
	defaultAccount = "sudo@admin"
)

// Here amount cannot be negative, bcoz 2 parties are involved
// payer must always pay something to payee
// instead of using negative amount switch payer and payee
func New(payee, payer string, amount float64, refId string) (*TransactHistory, error) {
	// sanity checks
	if len(payer) == 0 || len(payee) == 0 {
		return nil, fmt.Errorf("id of payer or payee cannot be empty")
	}

	if amount <= 0 {
		return nil, fmt.Errorf("amount cannot be zero or negative")
	}

	Account.UpdateBalanceForce(payer, -1*amount)
	Account.UpdateBalanceForce(payee, amount)

	transact := TransactHistory{
		id:     getId(payer, false, uint32(getUserTxsCount(payer))),
		amount: amount,
		payee:  payee,
		rev:    false,
		payer:  payer,
		refId:  refId,
	}

	transact.desc = transact.getDesc()

	addTx(&transact)

	// return by value isn't
	// necessary because user can't
	// modify the actual tx
	return &transact, nil
}

func NewDebitTx(userId string, amount float64, refId string) (*TransactHistory, error) {
	// sanity checks
	if amount <= 0 {
		return nil, fmt.Errorf("amount cannot be zero or negative")
	}
	return newIndividualTx(userId, -1*amount, refId)
}

func NewCreditTx(userId string, amount float64, refId string) (*TransactHistory, error) {
	// sanity checks
	if amount <= 0 {
		return nil, fmt.Errorf("amount cannot be zero or negative")
	}
	return newIndividualTx(userId, amount, refId)
}

// send new slices of pointers
// because user can't modify the actual tx
func GetTxs(payerId string) []*TransactHistory {
	transacts := getTxsOfPayer(payerId)
	newSlice := make([]*TransactHistory, len(transacts))
	for i, v := range transacts {
		newSlice[i] = &v
	}
	return newSlice
}

func GetAccountingTxs(userId string) []*TransactHistory {
	revTxs := getTxsOfPayee(userId)
	txs := getTxsOfPayer(userId)

	lenTxs := len(txs)
	lenRevTxs := len(revTxs)

	newSlice := make([]*TransactHistory, lenTxs+lenRevTxs)

	// copy all original txs
	for i, v := range txs {
		newSlice[i] = &v
	}

	// copy all reverse txs
	for i, v := range revTxs {
		newSlice[i+lenTxs] = &v
	}

	return newSlice
}

// helper internal funcs

func getId(userId string, rev bool, num uint32) string {
	if rev {
		return fmt.Sprint(userId, revIdDivider, num)
	}
	return fmt.Sprint(userId, idDivider, num)
}

func newIndividualTx(userId string, amount float64, refId string) (*TransactHistory, error) {
	// sanity checks
	if len(userId) == 0 {
		return nil, fmt.Errorf("userId cannot be empty")
	}

	Account.UpdateBalanceForce(userId, -1*amount)

	tx := TransactHistory{
		id:     getId(userId, false, uint32(getUserTxsCount(userId))),
		amount: amount,
		payee:  defaultAccount,
		rev:    false,
		payer:  userId,
		refId:  refId,
	}

	// update payer map
	addTx(&tx)

	return &tx, nil
}

func (t *TransactHistory) getDesc() string {
	if t.rev {
		return fmt.Sprint(t.payee, " was paid ", t.amount, " by ", t.payer)
	} else if len(t.payee) != 0 && len(t.payer) != 0 {
		return fmt.Sprint(t.payer, " paid ", t.amount, " to ", t.payee)
	} else if len(t.payer) != 0 {
		return fmt.Sprint(t.payer, " expensed ", t.amount)
	} else if len(t.payee) != 0 {
		return fmt.Sprint(t.payee, " received ", t.amount)
	}
	return fmt.Sprint("invalid transaction of value ", t.amount)
}

func addTx(tx *TransactHistory) {
	// create reverse for accounting
	revTransact := TransactHistory{
		id:     getId(tx.payee, true, uint32(getRevTxsCount(tx.payee))),
		amount: tx.amount,
		payer:  tx.payer,
		payee:  tx.payee,
		rev:    true,
		refId:  tx.refId,
	}

	revTransact.desc = revTransact.getDesc()

	// update payer map
	addActualTx(tx)

	// update payee map
	addAccountingTx(&revTransact)
}
