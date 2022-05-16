package core

import "fmt"

func GetPayablesFromMap(m map[string]float64) ([]*Payable, error) {
	bal := make([]*BalanceAccount, len(m))
	i := 0
	for k, v := range m {
		bal[i] = &BalanceAccount{k, v}
		i++
	}
	return GetPayableFromBalances(bal)
}

func GetPayableFromBalances(bal []*BalanceAccount) ([]*Payable, error) {
	ok, err := sanitizeAndCheckBalances(bal)
	if !ok {
		return nil, err
	}
	refMap := map[string]float64{}
	for _, v := range bal {
		if v.Amount != 0 {
			refMap[v.UserId] = v.Amount
		}
	}
	payables := make([]*Payable, 0)
	//TODO: main logic here
	return payables, nil
}

func MinimizePayables(payables []*Payable) ([]*Payable, error) {
	// Convert payables to map
	m := map[string]float64{}
	for _, p := range payables {
		m[p.PayerId] -= p.Amount
		m[p.PayeeId] += p.Amount
	}
	return GetPayablesFromMap(m)
}

func GetBalancesFromMap(m map[string]float64) []BalanceAccount {
	payables := make([]BalanceAccount, 0)
	for k, v := range m {
		if v != 0 {
			payables = append(payables, BalanceAccount{k, v})
		}
	}
	return payables
}

func sanitizeAndCheckBalances(bal []*BalanceAccount) (bool, error) {
	total := 0.0
	for _, v := range bal {
		total += v.Amount
	}
	if total != 0 {
		return false, fmt.Errorf("total balance must be zero found %f", total)
	}
	return true, nil
}
