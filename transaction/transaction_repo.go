package transaction

var (
	txMap    = make(map[string][]TransactHistory)
	revTxMap = make(map[string][]TransactHistory)
)

var getTxsOfPayer = func(payerId string) []TransactHistory {
	val, ok := txMap[payerId]
	if !ok {
		val = make([]TransactHistory, 0)
		txMap[payerId] = val
	}
	return val
}

var getTxsOfPayee = func(payeeId string) []TransactHistory {
	val, ok := revTxMap[payeeId]
	if !ok {
		val = make([]TransactHistory, 0)
		revTxMap[payeeId] = val
	}
	return val
}

var addActualTx = func(tx *TransactHistory) {
	txs := getTxsOfPayer(tx.payer)
	txs = append(txs, *tx)
	txMap[tx.payer] = txs
}

var addAccountingTx = func(tx *TransactHistory) {
	txs := getTxsOfPayee(tx.payee)
	txs = append(txs, *tx)
	revTxMap[tx.payee] = txs
}

var getUserTxsCount = func(userId string) int {
	return len(getTxsOfPayer(userId))
}

var getRevTxsCount = func(userId string) int {
	return len(getTxsOfPayee(userId))
}
