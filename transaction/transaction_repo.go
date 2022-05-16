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

var addPayerTx = func(payerId string, tx *TransactHistory) {
	txs := getTxsOfPayer(payerId)
	txs = append(txs, *tx)
	txMap[payerId] = txs
}

var addPayeeTx = func(payeeId string, tx *TransactHistory) {
	txs := getTxsOfPayee(payeeId)
	txs = append(txs, *tx)
	revTxMap[payeeId] = txs
}

var getUserTxsCount = func(userId string) int {
	return len(getTxsOfPayer(userId))
}

var getRevTxsCount = func(userId string) int {
	return len(getTxsOfPayee(userId))
}
