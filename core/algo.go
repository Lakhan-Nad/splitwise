package core

import (
	"fmt"
)

// Public API

func MinimizePayables(payables []*Payable) ([]*Payable, error) {
	// Convert payables to map
	m := map[string]float64{}
	for _, p := range payables {
		if p.Amount < 0 {
			return nil, fmt.Errorf("amount must be positive")
		}
		m[p.PayerId] -= p.Amount
		m[p.PayeeId] += p.Amount
	}
	return GetPayablesFromMap(m)
}

func GetPayablesFromBalanceAccounts(bal []*BalanceAccount) ([]*Payable, error) {
	balances := make([]BalanceAccount, 0)
	for _, b := range bal {
		if b.UserId == "" {
			return nil, fmt.Errorf("userId is empty")
		}
		if b.Amount != 0 {
			balances = append(balances, *b)
		}
	}
	return UseBalanceMapSimplifyToPayableStrategy(balances, &CashFlowMinimiseStrategy{})
}

func GetPayablesFromMap(m map[string]float64) ([]*Payable, error) {
	balances := make([]BalanceAccount, 0)
	for k, v := range m {
		if k == "" {
			return nil, fmt.Errorf("userId is empty")
		}
		if v != 0 {
			balances = append(balances, BalanceAccount{
				UserId: k,
				Amount: v,
			})
		}
	}
	return UseBalanceMapSimplifyToPayableStrategy(balances, &CashFlowMinimiseStrategy{})
}
