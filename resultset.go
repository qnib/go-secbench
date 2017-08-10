package secbench

import (
	"sync"
)

type ResultSets struct {
	mu sync.Mutex
	data map[string]Rule

}


func NewResultSet() *ResultSets {
	rs := &ResultSets{
		data: map[string]Rule{},
	}
	return rs
}

func (rs *ResultSets) AddRule(num, desc, mode string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.data[num] = NewRule(num, desc, mode)
}

func (rs *ResultSets) AddResult(num, mode, msg string) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	res := NewResult(mode, msg)
	rs.data[num].AddResult(res)
}
