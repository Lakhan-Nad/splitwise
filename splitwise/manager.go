package splitwise

import (
	"fmt"

	Core "workat.tech/splitwise/core"
	Splits "workat.tech/splitwise/splits"
	Transation "workat.tech/splitwise/transaction"
)

type SplitExpense struct {
	contributions []*Contribution
	amount        float64
	strategy      Splits.SplitStrategy
	splits        []*Splits.Split
}

func (se *SplitExpense) Validate() (bool, error) {
	if se.amount < 0 {
		return false, fmt.Errorf("cannot split negative amount")
	}

	ok, err := se.strategy.ValidateSplits(se.amount, se.splits)
	if !ok {
		return false, err
	}

	return true, nil
}

//TODO: use better refId instead of xyz
func (se *SplitExpense) Calculate() (bool, error) {
	// eager checking of valid split expense given or not
	validated, validationErr := se.Validate()
	if !validated {
		return false, validationErr
	}

	// calculate split shares
	Splits.UseSplitStrategy(se.strategy, se.amount, se.splits)

	// make map for balances
	m := make(map[string]float64)

	// add expense for each contributor
	// as a credit from their account
	for _, e := range se.contributions {
		Transation.NewCreditTx(e.UserId(), se.amount, "xyz")
		m[e.UserId()] = m[e.UserId()] + e.Amount()
	}

	// add individual user's transaction to payer
	for _, s := range se.splits {
		m[s.UserId()] = m[s.UserId()] - float64(s.Share())
	}

	payables, err := Core.GetPayablesFromMap(m)

	if err != nil {
		return false, err
	}

	// add payables as transactions
	for _, p := range payables {
		Transation.New(p.PayeeId, p.PayerId, p.Amount, "xyz")
	}

	return true, nil
}

// contribution amount can never be negative or zero
func New(contributions []*Contribution, splits []*Splits.Split, strategy Splits.SplitStrategy) *SplitExpense {
	totalAmount := 0.0
	for _, e := range contributions {
		totalAmount += e.Amount()
	}
	return &SplitExpense{
		contributions: contributions,
		splits:        splits,
		strategy:      strategy,
		amount:        totalAmount,
	}
}

//Notes

// 1. Use RefId as Id of expense to store ExpenseMetaData
// 2. Use RefId of each split as Id of each split to store SplitMetaData
// 3. Balance of User is sum of all it's Accounting txs
// 4. For simpliying the split as mentioned in 4h optional requirement
//    use balance all actual transactions and create a simplified
//    graph of outstanding amounts between each user (exclude edges with 0 amount)
//    this will give us a DAG (Directed Acyclic Graph)
//    then use topological sort in graph to find simplified outstanding
