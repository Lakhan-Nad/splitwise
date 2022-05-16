package core

import (
	"container/heap"
	"fmt"
	"math"
)

// Public API

func MinimizePayables(payables []*Payable) ([]*Payable, error) {
	// Convert payables to map
	m := map[string]float64{}
	for _, p := range payables {
		m[p.PayerId] -= p.Amount
		m[p.PayeeId] += p.Amount
	}
	return GetPayablesFromMap(m)
}

func GetPayablesFromBalanceAccounts(bal []*BalanceAccount) ([]*Payable, error) {
	refMap := map[string]float64{}
	for _, b := range bal {
		if b.Amount != 0 && b.UserId != "" {
			refMap[b.UserId] = b.Amount
		} else if b.UserId == "" {
			return nil, fmt.Errorf("userId is empty")
		}
	}
	return getPayables(refMap)
}

func GetPayablesFromMap(m map[string]float64) ([]*Payable, error) {
	refMap := make(map[string]float64, len(m))
	for k, v := range m {
		if v != 0 && k != "" {
			refMap[k] = v
		} else if k == "" {
			return nil, fmt.Errorf("userId is empty")
		}
	}
	return getPayables(refMap)
}

// Private Implementation

func getPayables(refMap map[string]float64) ([]*Payable, error) {
	ok, err := sanitizeAndCheckBalances(refMap)
	if !ok {
		return nil, err
	}
	usersPayer := accountsHeap{
		s: []string{},
		m: refMap,
	}
	usersPayee := accountsHeap{
		s: []string{},
		m: refMap,
	}
	payables := make([]*Payable, 0)
	for k, v := range refMap {
		if v < 0 {
			usersPayer.s = append(usersPayer.s, k)
		} else {
			usersPayee.s = append(usersPayee.s, k)
		}
	}
	heap.Init(&usersPayer)
	heap.Init(&usersPayee)
	for usersPayer.Len() > 0 && usersPayee.Len() > 0 {
		p := &Payable{}
		p.PayerId = heap.Pop(&usersPayer).(string)
		p.PayeeId = heap.Pop(&usersPayee).(string)
		p.Amount = math.Max(math.Abs(refMap[p.PayerId]), refMap[p.PayeeId])
		refMap[p.PayerId] += p.Amount
		refMap[p.PayeeId] -= p.Amount
		if refMap[p.PayerId] != 0 {
			heap.Push(&usersPayer, p.PayerId)
		}
		if refMap[p.PayeeId] != 0 {
			heap.Push(&usersPayee, p.PayeeId)
		}
		payables = append(payables, p)
	}
	return payables, nil
}

func sanitizeAndCheckBalances(bal map[string]float64) (bool, error) {
	total := 0.0
	for _, v := range bal {
		total += v
	}
	if total != 0 {
		return false, fmt.Errorf("total balance must be zero found %f", total)
	}
	return true, nil
}

// heap implementation
type accountsHeap struct {
	s []string
	m map[string]float64
}

func (a *accountsHeap) Len() int {
	return len(a.s)
}

func (a *accountsHeap) Less(i, j int) bool {
	return math.Abs(a.m[a.s[i]]) < math.Abs(a.m[a.s[j]])
}

func (a *accountsHeap) Swap(i, j int) {
	a.s[i], a.s[j] = a.s[j], a.s[i]
}

func (a *accountsHeap) Push(x interface{}) {
	a.s = append(a.s, x.(string))
}

func (a *accountsHeap) Pop() interface{} {
	n := len(a.s)
	x := a.s[n-1]
	a.s = a.s[:n-1]
	return x
}
