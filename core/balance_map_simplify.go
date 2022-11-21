package core

import (
	"fmt"
)

type BalanceAccountSimplifyStrategy interface {
	Simplify(balances []BalanceAccount) ([]*Payable, error)
}

func UseBalanceMapSimplifyToPayableStrategy(bal []BalanceAccount, strategy BalanceAccountSimplifyStrategy) ([]*Payable, error) {
	// Sanitize and check balances
	ok, err := sanitizeAndCheckBalances(bal)
	if !ok {
		return nil, err
	}
	return strategy.Simplify(bal)
}

func sanitizeAndCheckBalances(bal []BalanceAccount) (bool, error) {
	total := 0.0
	for _, v := range bal {
		total += v.Amount
		if v.AccountId == "" {
			return false, fmt.Errorf("payer  is empty")
		}
	}
	if total != 0 {
		return false, fmt.Errorf("total balance must be zero found %f", total)
	}
	return true, nil
}
