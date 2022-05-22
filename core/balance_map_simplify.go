package core

import (
	"container/heap"
	"fmt"
	"math"
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
		if v.UserId == "" {
			return false, fmt.Errorf("userId is empty")
		}
	}
	if total != 0 {
		return false, fmt.Errorf("total balance must be zero found %f", total)
	}
	return true, nil
}

// heap startegy for simplify payables
// How it works:
// 1. it first segreggates the balances into two groups:
//		a. positive balances
//		b. negative balances
// 2. it always tries to first pair out same amount of balances from the two groups
// 3. if unable to find a pair
// 		a. it tries to make a edge between maximum absolute balance from both group,
//          uses heap for the implementation of same
//      b. if a balance has any amount left it pushes that back to heap

type CashFlowMinimiseStrategy struct {
}

type accountsSet map[*BalanceAccount]bool
type amountToAccountsMap map[float64]accountsSet

func (strategy *CashFlowMinimiseStrategy) Simplify(balances []BalanceAccount) ([]*Payable, error) {
	posMap := make(amountToAccountsMap)
	negMap := make(amountToAccountsMap)
	negHeap := accountsHeap{
		s: make([]*BalanceAccount, 0),
	}
	posHeap := accountsHeap{
		s: make([]*BalanceAccount, 0),
	}
	heap.Init(&posHeap)
	heap.Init(&negHeap)
	payables := make([]*Payable, 0)
	for _, v := range balances {
		if v.Amount > 0 {
			_, ok := posMap[v.Amount]
			if !ok {
				posMap[v.Amount] = make(accountsSet)
			}
			posMap[v.Amount][&v] = true
		}
	}
	for _, v := range balances {
		if v.Amount < 0 {
			account := removeAndGetFirstBalance(posMap, -v.Amount)
			if account != nil {
				payables = append(payables, newPayableFromBalance(&v, account))
				continue
			}
			negMap[v.Amount] = make(accountsSet)
			negMap[v.Amount][&v] = true
		}
	}
	for _, v := range balances {
		if v.Amount < 0 {
			heap.Push(&negHeap, &v)
		} else {
			heap.Push(&posHeap, &v)
		}
	}
	for negHeap.Len() > 0 && posHeap.Len() > 0 {
		payer := negHeap.Peek()
		if payer.Amount == 0 {
			continue
		}
		payee := posHeap.Peek()
		if payee.Amount == 0 {
			continue
		}
		heap.Pop(&negHeap)
		heap.Pop(&posHeap)
		payable := newPayableFromBalance(payer, payee)
		payables = append(payables, payable)
		account := removeAndGetFirstBalance(posMap, -payer.Amount)
		if account != nil {
			payables = append(payables, newPayableFromBalance(payer, account))
		}
		account = removeAndGetFirstBalance(negMap, -payee.Amount)
		if account != nil {
			payables = append(payables, newPayableFromBalance(account, payee))
		}
		if payer.Amount != 0 {
			heap.Push(&negHeap, payer)
		} else if payee.Amount != 0 {
			heap.Push(&posHeap, payee)
		}
	}
	return payables, nil
}

func newPayableFromBalance(payerBal *BalanceAccount, payeeBal *BalanceAccount) *Payable {
	p := &Payable{}
	p.Amount = math.Min(math.Abs(payerBal.Amount), payeeBal.Amount)
	p.PayeeId = payeeBal.UserId
	p.PayerId = payerBal.UserId
	payeeBal.Amount -= p.Amount
	payerBal.Amount += p.Amount
	return p
}

func removeAndGetFirstBalance(aMap amountToAccountsMap, amount float64) *BalanceAccount {
	m, ok := aMap[amount]
	if ok && len(m) > 0 {
		var oneAccount *BalanceAccount
		for k := range m {
			oneAccount = k
			break
		}
		delete(m, oneAccount)
		return oneAccount
	}
	return nil
}

// heap implementation
type accountsHeap struct {
	s []*BalanceAccount
}

func (a *accountsHeap) Len() int {
	return len(a.s)
}

func (a *accountsHeap) Less(i, j int) bool {
	// we need max heap
	return math.Abs(a.s[i].Amount) > math.Abs(a.s[j].Amount)
}

func (a *accountsHeap) Swap(i, j int) {
	a.s[i], a.s[j] = a.s[j], a.s[i]
}

func (a *accountsHeap) Push(x interface{}) {
	a.s = append(a.s, x.(*BalanceAccount))
}

func (a *accountsHeap) Pop() interface{} {
	n := len(a.s)
	x := a.s[n-1]
	a.s = a.s[:n-1]
	return x
}

func (a *accountsHeap) Peek() *BalanceAccount {
	if a.Len() > 0 {
		return a.s[0]
	}
	return nil
}
