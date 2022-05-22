package core

import "fmt"

type PayableSimplifyStrategy interface {
	Simplify(payables []*Payable) ([]*Payable, error)
}

func UsePaybleSimplifyStrategy(payables []*Payable, strategy PayableSimplifyStrategy) ([]*Payable, error) {
	ok, err := sanitizeAndCheckPayables(payables)
	if !ok {
		return nil, err
	}
	return strategy.Simplify(payables)
}

func sanitizeAndCheckPayables(payables []*Payable) (bool, error) {
	for _, p := range payables {
		if p.Amount < 0 {
			return false, fmt.Errorf("amount must be positive")
		}
		if p.PayeeId == "" || p.PayerId == "" {
			return false, fmt.Errorf("userId is empty")
		}
	}
	return true, nil
}
