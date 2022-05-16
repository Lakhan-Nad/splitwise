package account

import (
	"fmt"
	"time"
)

var (
	accountMap = make(map[string]*TransactionAccount)
)

func NoAccountError(userId string) error {
	return fmt.Errorf("user %s does not have a transaction account", userId)
}

func AccountExistError(userId string) error {
	return fmt.Errorf("user %s already has a transaction account", userId)
}

var addNewAccount = func(userId string, openingBalance float64) {
	accountMap[userId] = &TransactionAccount{
		userId:      userId,
		balance:     openingBalance,
		inUse:       true,
		lastUpdated: time.Now().Unix(),
	}
}

var upsertTransactionAccount = func(userId string, balance float64) *TransactionAccount {
	if val, ok := accountMap[userId]; !ok {
		return val
	}
	addNewAccount(userId, balance)
	return accountMap[userId]
}

var AddNewTransactionAccount = func(userId string, openingBalance float64) error {
	if _, ok := accountMap[userId]; !ok {
		addNewAccount(userId, openingBalance)
		return nil
	}
	return AccountExistError(userId)
}
var UpdateBalanceForce = func(userId string, balance float64) {
	val, ok := accountMap[userId]
	if !ok {
		val = upsertTransactionAccount(userId, balance)
	}
	val.balance += balance
	val.lastUpdated = time.Now().Unix()
}

var UpdateBalance = func(userId string, balance float64) error {
	val, ok := accountMap[userId]
	if !ok {
		return NoAccountError(userId)
	}
	val.balance += balance
	val.lastUpdated = time.Now().Unix()
	return nil
}

var GetBalanceForce = func(userId string) float64 {
	val, ok := accountMap[userId]
	if !ok {
		val = upsertTransactionAccount(userId, 0)
	}
	return val.balance
}

var GetBalance = func(userId string) (float64, error) {
	val, ok := accountMap[userId]
	if !ok {
		return 0, NoAccountError(userId)
	}
	return val.balance, nil
}

var GetLastUpdatedForce = func(userId string) int64 {
	val, ok := accountMap[userId]
	if !ok {
		val = upsertTransactionAccount(userId, 0)
	}
	return val.lastUpdated
}

var GetLastUpdated = func(userId string) (int64, error) {
	val, ok := accountMap[userId]
	if !ok {
		return 0, NoAccountError(userId)
	}
	return val.lastUpdated, nil
}
